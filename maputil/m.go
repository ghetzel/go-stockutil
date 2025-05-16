package maputil

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	utilutil "github.com/ghetzel/go-stockutil/utils"
)

var MapXmlRootTagName = `data`
var MapXmlStructTagName = `xml`

type MapSetFunc func(m *Map, key string) any

type IterOptions struct {
	TagName  string
	SortKeys bool
}

type Item struct {
	Key   any
	Value any
	K     string
	V     typeutil.Variant
	m     *Map
}

func (self *Item) Set(value any) error {
	if self.m == nil {
		return fmt.Errorf("cannot set value: no parent Map")
	} else {
		var nv = self.m.Set(self.K, value)
		self.V = nv
		self.Value = self.V.Value
		return nil
	}
}

type ItemFunc func(key string, value typeutil.Variant) error
type KeyTransformFunc func(string) string

// A Map object (or "M" object) is a utility struct that makes it straightforward to
// work with interface data types that contain map-like data (has a reflect.Kind equal
// to reflect.Map).
type Map struct {
	data              any
	structTagKey      string
	rootTagName       string
	xmlMarshalGeneric bool
	xmlKeyTransformFn KeyTransformFunc
	atomic            *sync.Mutex
}

func NewMap() *Map {
	return M(nil)
}

// Create a new Variant map object from the given value.  A wide range of values are accepted, and
// the best effort is made to convert those values into a usable map. Accepted values include typeutil.Variant,
// any value with a reflect.Kind of reflect.Map, sync.Map, another maputil.Map, url.Values,
// http.Header, or a string or []byte which will be decoded using json.Unmarshal if and only if the
// string begins with "{" and ends with "}".
func M(data any) *Map {
	if dataV, ok := data.(typeutil.Variant); ok {
		data = dataV.Value
	} else if dataM, ok := data.(*Map); ok {
		return dataM
	} else if dataM, ok := data.(Map); ok {
		return &dataM
	} else if dataSM, ok := data.(*sync.Map); ok {
		var dataM = make(map[string]any)
		dataSM.Range(func(key, value any) bool {
			dataM[typeutil.String(key)] = value
			return true
		})

		data = dataM
	} else if uV, ok := data.(url.Values); ok {
		var dataM = make(map[string]any)

		for k, v := range uV {
			switch len(v) {
			case 0:
				break
			case 1:
				dataM[k] = typeutil.Auto(v[0])
			default:
				dataM[k] = sliceutil.Autotype(v)
			}
		}

		data = dataM
	} else if hV, ok := data.(http.Header); ok {
		var dataM = make(map[string]any)

		for k, v := range hV {
			switch len(v) {
			case 0:
				break
			case 1:
				dataM[k] = typeutil.Auto(v[0])
			default:
				dataM[k] = sliceutil.Autotype(v)
			}
		}

		data = dataM
	} else if dS, ok := data.(string); ok {
		if stringutil.IsSurroundedBy(strings.TrimSpace(dS), `{`, `}`) {
			data = make(map[string]any)
			json.Unmarshal([]byte(dS), &data)
		}
	} else if dB, ok := data.([]byte); ok {
		if stringutil.IsSurroundedBy(strings.TrimSpace(string(dB)), `{`, `}`) {
			data = make(map[string]any)
			json.Unmarshal(dB, &data)
		}
	}

	if data == nil {
		data = make(map[string]any)
	}

	return &Map{
		data:         data,
		structTagKey: UnmarshalStructTag,
	}
}

// Specify which struct tag to honor for generating field names when then
// underlying data is a struct.
func (self *Map) Tag(key string) *Map {
	self.structTagKey = key
	return self
}

// Return the underlying value the M-object was created with.
func (self *Map) Value() any {
	return self.data
}

// Delete a value from the map.
func (self *Map) Delete(key string) {
	Delete(self.data, key)
}

// internal: acquire write lock
func (self *Map) lock() {
	if self.atomic == nil {
		self.atomic = new(sync.Mutex)
	}

	self.atomic.Lock()
}

// internal: release write lock
func (self *Map) unlock() {
	if self.atomic != nil {
		self.atomic.Unlock()
	}
}

// internal: unlocked implementation of set()
func (self *Map) set(key string, value any) typeutil.Variant {
	var vv = typeutil.V(value)
	self.data = DeepSet(self.data, strings.Split(key, `.`), vv)

	return vv
}

// Set a value in the Map at the given dot.separated key to a value.
func (self *Map) Set(key string, value any) typeutil.Variant {
	self.lock()
	defer self.unlock()

	return self.set(key, value)
}

// Set a value in the Map using a function.  The map will be locked to
// other modifications for the duration of the function's execution.
func (self *Map) SetFunc(key string, vfunc MapSetFunc) typeutil.Variant {
	if vfunc != nil {
		self.lock()
		defer self.unlock()

		return self.set(key, vfunc(self, key))
	}

	return typeutil.V(nil)
}

// Set a value in the Map at the given dot.separated key to a value, but only if the
// current value at that key is that type's zero value.
func (self *Map) SetIfZero(key string, value any) (typeutil.Variant, bool) {
	self.lock()
	defer self.unlock()

	if v := self.Get(key); v.IsZero() {
		return self.set(key, value), true
	} else {
		return v, false
	}
}

// Set a value in the Map at the given dot.separated key to a value, but only if the
// new value is not a zero value.
func (self *Map) SetValueIfNonZero(key string, value any) (typeutil.Variant, bool) {
	self.lock()
	defer self.unlock()

	if !typeutil.IsZero(value) {
		return self.set(key, value), true
	} else {
		return self.Get(key), false
	}
}

// Copy the items from a map into this one.
func (self *Map) Merge(other any) int {
	self.lock()
	defer self.unlock()

	var count int

	for k, v := range M(other).MapNative() {
		self.set(k, v)
		count += 1
	}

	return count
}

// Reject all nil values from the map.
func (self *Map) Compact() *Map {
	self.lock()
	defer self.unlock()

	var d = self.MapNative()

	for k, v := range d {
		if v == nil {
			delete(d, k)
		}
	}

	self.data = d
	return self
}

// Performs a JSONPath query against the given object and returns the results.
// See JSONPath for details.
func (self *Map) JSONPath(query string, fallback ...any) any {
	if d, err := JSONPath(self.MapNative(), query); err == nil && d != nil {
		return d
	}

	return sliceutil.First(fallback)
}

// Retrieve a value from the Map by the given dot.separated key, or return a fallback
// value.  Return values are a typeutil.Variant, which can be easily coerced into
// various types.
func (self *Map) Get(key string, fallbacks ...any) typeutil.Variant {
	var native = self.MapNative(self.structTagKey)

	if v, ok := native[key]; ok && v != nil {
		return typeutil.Variant{
			Value: v,
		}
	} else if v := DeepGet(self.data, strings.Split(key, `.`), fallbacks...); v != nil {
		return typeutil.Variant{
			Value: v,
		}
	} else {
		return typeutil.Variant{}
	}
}

// Return the value at key as an automatically converted value.
func (self *Map) Auto(key string, fallbacks ...any) any {
	return self.Get(key, fallbacks...).Auto()
}

// Return the value at key as a string.
func (self *Map) String(key string, fallbacks ...any) string {
	return self.Get(key, fallbacks...).String()
}

// Return the value at key interpreted as a Time.
func (self *Map) Time(key string, fallbacks ...any) time.Time {
	return self.Get(key, fallbacks...).Time()
}

// Return the value at key interpreted as a Duration.
func (self *Map) Duration(key string, fallbacks ...any) time.Duration {
	return self.Get(key, fallbacks...).Duration()
}

// Return the value at key as a bool.
func (self *Map) Bool(key string) bool {
	return self.Get(key).Bool()
}

// Return the value at key as an integer.
func (self *Map) Int(key string, fallbacks ...any) int64 {
	return self.Get(key, fallbacks...).Int()
}

// Return the value at key as a native integer.
func (self *Map) NInt(key string, fallbacks ...any) int {
	return self.Get(key, fallbacks...).NInt()
}

// Return the value at key as a float.
func (self *Map) Float(key string, fallbacks ...any) float64 {
	return self.Get(key, fallbacks...).Float()
}

// Return the value at key as a byte slice.
func (self *Map) Bytes(key string) []byte {
	return self.Get(key).Bytes()
}

// Return the value at key as a slice.  Scalar values will be returned as a slice containing
// only that value.
func (self *Map) Slice(key string) []typeutil.Variant {
	return self.Get(key).Slice()
}

// Convert the given value to a slice using typeutil.Slice, then return each element as a Map.
func (self *Map) SliceOfMaps(key string) []*Map {
	return SliceOfMaps(self.Get(key).Value)
}

// Same as Slice(), but returns a []string
func (self *Map) Strings(key string) []string {
	return self.Get(key).Strings()
}

// Return the value at key as an error, or nil if the value is not an error.
func (self *Map) Err(key string) error {
	return self.Get(key).Err()
}

// Return the value at key as a Map.  If the resulting value is nil or not a
// map type, a null Map will be returned.  All values retrieved from a null
// Map will return that type's zero value.
func (self *Map) Map(key string, tagName ...string) map[typeutil.Variant]typeutil.Variant {
	if len(tagName) == 0 {
		tagName = []string{self.structTagKey}
	}

	return self.Get(key).Map(tagName...)
}

// Return the value as a map[string]any.
func (self *Map) MapNative(tagName ...string) map[string]any {
	if len(tagName) == 0 {
		tagName = []string{self.structTagKey}
	}

	return typeutil.MapNative(self.data, tagName...)
}

// Return the value as a map[string]string.
func (self *Map) MapString(tagName ...string) map[string]string {
	var rv = make(map[string]string)

	for k, v := range self.MapNative(tagName...) {
		rv[k] = typeutil.String(v)
	}

	return rv
}

func (self *Map) JSON(indent ...string) (data []byte) {
	if len(indent) > 0 {
		data, _ = json.MarshalIndent(self.data, ``, indent[0])
	} else {
		data, _ = json.Marshal(self.data)
	}

	return
}

func (self *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.data)
}

func (self *Map) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &self.data)
}

// Uses the extended Sprintf in this package, passing this map as the data used in the given format string.
func (self *Map) Sprintf(format string) string {
	return Sprintf(format, self.MapNative())
}

// Uses the extended Fprintf in this package, passing this map as the data used in the given format string.
func (self *Map) Fprintf(w io.Writer, format string) {
	Fprintf(w, format, self.MapNative())
}

func xn(generic bool, ifGeneric string, otherwise string) xml.Name {
	if generic {
		return _xn(ifGeneric)
	} else {
		return _xn(otherwise)
	}
}

func _xn(tagName string) xml.Name {
	return xml.Name{
		Local: tagName,
	}
}

func xt(value any) string {
	if typeutil.IsMap(value) {
		return `object`
	} else if typeutil.IsArray(value) {
		return `array`
	} else {
		return strings.TrimPrefix(stringutil.Hyphenate(fmt.Sprintf("%T", value)), `-`)
	}
}

func (self *Map) valueToXmlTokens(parent *xml.StartElement, value any, key string) (tokens []xml.Token, ferr error) {
	var g = self.xmlMarshalGeneric

	if self.xmlKeyTransformFn != nil {
		key = self.xmlKeyTransformFn(key)
	}

	if typeutil.IsScalar(value) {
		var open = xml.StartElement{
			Name: xn(g, `item`, key),
		}

		if g {
			open.Attr = []xml.Attr{
				{
					Name:  _xn(`key`),
					Value: key,
				}, {
					Name:  _xn(`type`),
					Value: utilutil.DetectConvertType(value).String(),
				},
			}
		}

		tokens = append(tokens, open, xml.CharData(typeutil.String(value)), xml.EndElement{
			Name: xn(g, `item`, key),
		})

	} else if typeutil.IsArray(value) {
		var start = xml.StartElement{
			Name: xn(g, `item`, key),
		}

		if g {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  _xn(`type`),
				Value: xt(value),
			}, xml.Attr{
				Name:  _xn(`key`),
				Value: key,
			})
		}

		tokens = append(tokens, start)

		for i, v := range sliceutil.Sliceify(value) {
			if ts, err := self.valueToXmlTokens(&start, v, `element`); err == nil {
				tokens = append(tokens, ts...)
			} else {
				ferr = fmt.Errorf("[%d]: %v", i, err)
				return
			}
		}

		tokens = append(tokens, xml.EndElement{
			Name: start.Name,
		})
	} else {
		var children = M(value).MapNative(MapXmlStructTagName)
		var ckeys = StringKeys(children)
		sort.Strings(ckeys)

		var start = xml.StartElement{
			Name: xn(g, `item`, key),
		}

		if g {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  _xn(`type`),
				Value: xt(value),
			}, xml.Attr{
				Name:  _xn(`key`),
				Value: key,
			})
		}

		tokens = append(tokens, start)

		for _, k := range ckeys {
			var v = children[k]

			if ts, err := self.valueToXmlTokens(&start, v, k); err == nil {
				tokens = append(tokens, ts...)
			} else {
				ferr = fmt.Errorf("%s: %v", k, err)
				return
			}
		}

		tokens = append(tokens, xml.EndElement{
			Name: start.Name,
		})
	}

	return
}

func (self *Map) SetMarshalXmlGeneric(yes bool) {
	self.xmlMarshalGeneric = yes
}

// set the name of the root XML tag, used by MarshalXML.
func (self *Map) SetRootTagName(root string) {
	self.rootTagName = root
}

// Set a function that will be used to generate XML tag names when calling MarshalXML.  This works
// for all keys, including ones that appear inside of maps.
func (self *Map) SetMarshalXmlKeyFunc(fn KeyTransformFunc) {
	self.xmlKeyTransformFn = fn
}

// Marshals the current data into XML.  Nested maps are output as nested elements.  Map values that
// are scalars (strings, numbers, bools, dates/times) will appear as attributes on the parent element.
func (self *Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var root = MapXmlRootTagName

	if self.rootTagName != `` {
		root = self.rootTagName
	}

	if self.xmlKeyTransformFn != nil {
		root = self.xmlKeyTransformFn(root)
	}

	start.Name = _xn(root)
	var tokens = []xml.Token{start}

	var children = self.MapNative(MapXmlStructTagName)
	var ckeys = StringKeys(children)
	sort.Strings(ckeys)

	for _, k := range ckeys {
		var v = children[k]

		if ts, err := self.valueToXmlTokens(&start, v, k); err == nil {
			tokens = append(tokens, ts...)
		} else {
			return err
		}
	}

	tokens = append(tokens, xml.EndElement{
		Name: start.Name,
	})

	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}

	return e.Flush()
}

// Return whether the value at the given key is that type's zero value.
func (self *Map) IsZero(key string) bool {
	return self.Get(key).IsZero()
}

// Return the keys in this Map object.  You may specify the name of a struct tag on the underlying
// object to use for generating key names.
func (self *Map) Keys(tagName ...string) []any {
	return Keys(self.MapNative(tagName...))
}

// A string slice version of Keys()
func (self *Map) StringKeys(tagName ...string) []string {
	return sliceutil.Stringify(self.Keys(tagName...))
}

// Return the length of the Map.
func (self *Map) Len() int {
	return len(self.MapNative())
}

// Iterate through each item in the map.
func (self *Map) Each(fn ItemFunc, tagName ...string) error {
	var tn = ``

	if len(tagName) > 0 && tagName[0] != `` {
		tn = tagName[0]
	}

	return self.each(fn, IterOptions{
		TagName: tn,
	})
}

// Iterate through each item in the map.
func (self *Map) each(fn ItemFunc, opts IterOptions) error {
	if fn != nil {
		var tn []string

		if opts.TagName != `` {
			tn = append(tn, opts.TagName)
		}

		var keys = self.StringKeys(tn...)

		if opts.SortKeys {
			sort.Strings(keys)
		}

		for _, key := range keys {
			if err := fn(key, self.Get(key)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *Map) Iter(opts ...IterOptions) <-chan Item {
	var itemchan = make(chan Item)

	if len(opts) == 0 {
		opts = []IterOptions{IterOptions{}}
	}

	go func() {
		self.each(func(key string, value typeutil.Variant) error {
			itemchan <- Item{
				Key:   key,
				Value: value.Value,
				K:     key,
				V:     value,
				m:     self,
			}

			return nil
		}, opts[0])

		close(itemchan)
	}()

	return itemchan
}

// A recursive walk form of Each()
func (self *Map) Walk(fn WalkFunc) error {
	return WalkStruct(self.data, fn)
}
