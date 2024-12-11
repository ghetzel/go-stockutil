// Utilities for converting, manipulating, and iterating over maps
package maputil

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/go-stockutil/utils"
	"github.com/mitchellh/mapstructure"
)

var rxJsonPathExpr = regexp.MustCompile(`\{.*?\}`)
var UnmarshalStructTag string = `maputil`
var SkipDescendants = errors.New("skip descendants")
var rxMapFmt = regexp.MustCompile(`(\$\{(?P<key>.*?)(?:\|(?P<fallback>.*?))?(?::(?P<fmt>%[^\}]+))?\})`) // ${key}, ${key:%04s}, ${key|fallback}

type WalkFunc func(value interface{}, path []string, isLeaf bool) error
type ApplyFunc func(key []string, value interface{}) (interface{}, bool)
type ConversionFunc func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error)
type deleteValue bool

type MergeOption int

const (
	AppendValues MergeOption = iota
)

type MergeOptions []MergeOption

func (self MergeOptions) Has(option MergeOption) bool {
	for _, opt := range self {
		if opt == option {
			return true
		}
	}

	return false
}

// Return an interface slice of the keys of the given map.
func Keys(input interface{}) []interface{} {
	var keys = make([]interface{}, 0)
	var rinput = typeutil.ResolveValue(input)

	if rinput == nil {
		return keys
	}

	var inputV = reflect.ValueOf(rinput)

	if inputV.Kind() == reflect.Map {
		keysV := inputV.MapKeys()

		for _, keyV := range keysV {
			keys = append(keys, keyV)
		}
	} else if syncMap, ok := input.(*sync.Map); ok {
		syncMap.Range(func(key interface{}, _ interface{}) bool {
			keys = append(keys, key)
			return true
		})
	}

	return keys
}

// Return a slice of strings representing the keys of the given map.
func StringKeys(input interface{}) []string {
	keys := sliceutil.Stringify(Keys(input))
	sort.Strings(keys)

	return keys
}

// Return the values from the given map.
func MapValues(input interface{}) []interface{} {
	var values = make([]interface{}, 0)
	var inputV = reflect.ValueOf(input)

	switch inputV.Kind() {
	case reflect.Map:
		for _, mapKeyV := range inputV.MapKeys() {
			if mapV := inputV.MapIndex(mapKeyV); mapV.IsValid() && mapV.CanInterface() {
				values = append(values, mapV.Interface())
			}
		}
	}

	return values
}

// Take an input map, and populate the struct instance pointed to by "populate".  Use the values of the tagname tag
// to inform which map keys should be used to fill struct fields, and if a Conversion function is given, that
// function will be used to allow values to be converted in preparation for becoming struct field values.
func TaggedStructFromMapFunc(input interface{}, populate interface{}, tagname string, converter ConversionFunc) error {
	if tagname == `` {
		tagname = UnmarshalStructTag
	}

	if converter == nil {
		converter = func(source reflect.Type, target reflect.Type, data interface{}) (interface{}, error) {
			if target.Kind() == reflect.String {
				return stringutil.ConvertToString(data)
			}

			if target.String() == `time.Time` || utils.IsTime(data) {
				return stringutil.ConvertToTime(data)
			}

			return data, nil
		}
	}

	if populateV, ok := populate.(reflect.Value); ok {
		if populateV.IsValid() && populateV.CanInterface() {
			populate = populateV.Interface()
		} else {
			return fmt.Errorf("Destination value is invalid or unsettable")
		}
	}

	var meta = new(mapstructure.Metadata)

	if decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           populate,
		TagName:          tagname,
		DecodeHook:       converter,
		WeaklyTypedInput: true,
		Metadata:         meta,
	}); err == nil {
		if err := decoder.Decode(input); err != nil {
			return fmt.Errorf("maputil: %v", err)
		}

		for _, field := range sliceutil.UniqueStrings(meta.Unused) {
			var key = strings.Split(field, `.`)
			var src = DeepGet(input, key)

			if utils.IsTime(src) {
				DeepSet(populate, key, typeutil.Time(src))

			} else if typeutil.IsMap(src) || typeutil.IsStruct(src) {
				for kv := range M(src).Iter() {
					DeepSet(populate, append(key, kv.K), kv.Value)
				}
			}
		}
	} else {
		return err
	}

	// fmt.Println(typeutil.Dump(populate))

	return nil
}

// Same as TaggedStructFromMapFunc, but does not perform any value conversion.
func TaggedStructFromMap(input interface{}, populate interface{}, tagname string) error {
	return TaggedStructFromMapFunc(input, populate, tagname, nil)
}

// Same as TaggedStructFromMapFunc, but no value conversion and uses the "maputil" struct tag.
func StructFromMap(input map[string]interface{}, populate interface{}) error {
	return TaggedStructFromMap(input, populate, ``)
}

// Join the given map, using innerJoiner to join keys and values, and outerJoiner to join the resulting key-value lines.
func Join(input interface{}, innerJoiner string, outerJoiner string) string {
	return DeepJoin(input, innerJoiner, outerJoiner, `.`)
}

// Join the given map, using innerJoiner to join keys and values, and outerJoiner to join the resulting key-value lines.
func DeepJoin(input interface{}, innerJoiner string, outerJoiner string, nestedSeparator string) string {
	parts := make([]string, 0)

	Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			parts = append(parts, strings.Join(path, nestedSeparator)+innerJoiner+stringutil.MustString(value))
		}

		return nil
	})

	return strings.Join(parts, outerJoiner)
}

// Split the given string, first on outerJoiner to form key-value lines, then each line on innerJoiner.
// Populates a map and returns the result.
func Split(input string, innerJoiner string, outerJoiner string) map[string]interface{} {
	rv := make(map[string]interface{})
	pairs := strings.Split(input, outerJoiner)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, innerJoiner, 2)

		if len(kv) == 2 {
			rv[kv[0]] = kv[1]
		}
	}

	return rv
}

// Take a flat (non-nested) map keyed with fields joined on fieldJoiner and return a
// deeply-nested map
func DiffuseMap(data map[string]interface{}, fieldJoiner string) (map[string]interface{}, error) {
	rv, _ := DiffuseMapTyped(data, fieldJoiner, "")
	return rv, nil
}

// Take a flat (non-nested) map keyed with fields joined on fieldJoiner and return a
// deeply-nested map
func DiffuseMapTyped(data map[string]interface{}, fieldJoiner string, typePrefixSeparator string) (map[string]interface{}, []error) {
	var errs = make([]error, 0)
	var output = make(map[string]interface{})

	//  get the list of keys and sort them because order in a map is undefined
	dataKeys := StringKeys(data)
	sort.Strings(dataKeys)

	//  for each data item
	for _, key := range dataKeys {
		var keyParts []string
		var value, _ = data[key]

		//  handle "typed" maps in which the type information is embedded
		if typePrefixSeparator != "" {
			var typeName string

			typeName, key = stringutil.SplitPairTrailing(key, typePrefixSeparator)

			if typeName == `` {
				typeName = `str`
			}

			if v, err := coerceIntoType(value, typeName); err == nil {
				value = v
			} else {
				errs = append(errs, err)
			}
		}

		keyParts = strings.Split(key, fieldJoiner)
		output = DeepSet(output, keyParts, value).(map[string]interface{})
	}

	return output, errs
}

// Take a deeply-nested map and return a flat (non-nested) map with keys whose intermediate tiers are joined with fieldJoiner
func CoalesceMap(data map[string]interface{}, fieldJoiner string) (map[string]interface{}, error) {
	return deepGetValues([]string{}, fieldJoiner, data), nil
}

// Take a deeply-nested map and return a flat (non-nested) map with keys whose intermediate tiers are joined with fieldJoiner
// Additionally, values will be converted to strings and keys will be prefixed with the datatype of the value
func CoalesceMapTyped(data map[string]interface{}, fieldJoiner string, typePrefixSeparator string) (map[string]interface{}, []error) {
	var errs = make([]error, 0)
	var rv = make(map[string]interface{})

	for k, v := range deepGetValues([]string{}, fieldJoiner, data) {
		if stringVal, err := stringutil.ToString(v); err == nil {
			rv[prepareCoalescedKey(k, v, typePrefixSeparator)] = stringVal
		} else {
			errs = append(errs, err)
		}
	}

	return rv, errs
}

func deepGetValues(keys []string, joiner string, data interface{}) map[string]interface{} {
	var rv = make(map[string]interface{})
	data = typeutil.ResolveValue(data)

	if data != nil {
		var dType = reflect.TypeOf(data)

		switch dType.Kind() {
		case reflect.Map:
			for k, v := range data.(map[string]interface{}) {
				var newKey = keys
				newKey = append(newKey, k)

				for kk, vv := range deepGetValues(newKey, joiner, v) {
					rv[kk] = vv
				}
			}

		case reflect.Slice, reflect.Array:
			for i, value := range sliceutil.Sliceify(data) {
				var newKey = keys
				newKey = append(newKey, strconv.Itoa(i))

				for k, v := range deepGetValues(newKey, joiner, value) {
					rv[k] = v
				}
			}

		default:
			rv[strings.Join(keys, joiner)] = data
		}
	}

	return rv
}

func prepareCoalescedKey(key string, value interface{}, typePrefixSeparator string) string {
	if typePrefixSeparator == "" {
		return key
	} else {
		var datatype string

		if dtype := utils.DetectConvertType(value); dtype != utils.Invalid {
			datatype = dtype.String()
		}

		return datatype + typePrefixSeparator + key
	}
}

func coerceIntoType(in interface{}, typeName string) (interface{}, error) {
	if dtype := stringutil.ParseType(typeName); dtype != stringutil.Invalid {
		if v, err := stringutil.ConvertTo(dtype, in); err == nil {
			return v, nil
		}
	}

	if inStr, err := stringutil.ToString(in); err == nil {
		return inStr, nil
	} else {
		return in, nil
	}
}

func Get(data interface{}, key string, fallback ...interface{}) interface{} {
	data = typeutil.ResolveValue(data)

	if typeutil.IsKind(data, reflect.Map) {
		var dataV = reflect.ValueOf(data)

		if valueV := dataV.MapIndex(reflect.ValueOf(key)); valueV.IsValid() {
			if valueI := valueV.Interface(); !typeutil.IsZero(valueI) {
				return valueI
			}
		}
	}

	if len(fallback) > 0 {
		return fallback[0]
	} else {
		return nil
	}
}

func DeepGet(data interface{}, path []string, fallbacks ...interface{}) interface{} {
	var current = typeutil.ResolveValue(data)

	if len(fallbacks) == 0 {
		fallbacks = []interface{}{nil}
	}

	var fallback = fallbacks[0]

	for i := 0; i < len(path); i++ {
		var part = path[i]
		var dValue = reflect.ValueOf(current)

		// if this value is not valid, return fallback here
		if !dValue.IsValid() {
			return fallback
		}

		var dType = dValue.Type()

		// for pointers and interfaces, get the underlying type
		switch dType.Kind() {
		case reflect.Interface, reflect.Ptr:
			dType = dType.Elem()
		}

		switch dType.Kind() {
		//  arrays
		case reflect.Slice, reflect.Array:
			if stringutil.IsInteger(part) {
				if partIndex, err := strconv.Atoi(part); err == nil {
					if partIndex < dValue.Len() {
						if value := dValue.Index(partIndex).Interface(); value != nil {
							current = value
							continue
						}
					}
				}
			} else if part == `*` {
				var subitems = make([]interface{}, dValue.Len())

				for j := 0; j < dValue.Len(); j++ {
					if value := dValue.Index(j).Interface(); value != nil {
						if i+1 < len(path) {
							subitems[j] = DeepGet(value, path[(i+1):], fallbacks...)
						} else {
							subitems[j] = value
						}
					} else {
						subitems[j] = fallback
					}
				}

				return subitems
			}

			return fallback

		//  maps
		case reflect.Map:
			if mapValue := dValue.MapIndex(
				reflect.ValueOf(part),
			); mapValue.IsValid() {
				current = mapValue.Interface()
			} else {
				return fallback
			}

		// structs
		case reflect.Struct:
			if dValue.Type().Kind() == reflect.Ptr {
				dValue = dValue.Elem()
			}

			if structField := dValue.FieldByName(part); structField.IsValid() && structField.CanInterface() {
				current = structField.Interface()
				continue
			}

		// attempting to retrieve nested data from a scalar value; return fallback
		default:
			return fallback
		}

	}

	return current
}

func DeepGetBool(data interface{}, path []string) bool {
	var vI = DeepGet(data, path, false)

	if v, ok := vI.(bool); ok && v {
		return true
	}

	return false
}

func DeepGetString(data interface{}, path []string) string {
	if v, err := stringutil.ToString(DeepGet(data, path)); err == nil {
		return v
	}

	return ``
}

// Delete a key to a given value in the given map.
func Delete(data interface{}, key interface{}) error {
	return Set(data, key, deleteValue(true))
}

// Set a key to a given value in the given map, reflect.Map Value, or slice/array.
func Set(data interface{}, key interface{}, value interface{}) error {
	var dataM reflect.Value
	var isDelete bool

	if _, ok := value.(deleteValue); ok {
		isDelete = true
	}

	if v, ok := data.(reflect.Value); ok {
		dataM = v
	} else {
		dataM = reflect.ValueOf(data)
	}

	// some shortcuts for common cases
	if asMap, ok := data.(map[string]interface{}); ok {
		if isDelete {
			delete(asMap, typeutil.String(key))
		} else {
			asMap[typeutil.String(key)] = value
		}

		return nil
	} else if dataM.CanInterface() {
		if asMap, ok := dataM.Interface().(map[string]interface{}); ok {
			if isDelete {
				delete(asMap, typeutil.String(key))
			} else {
				asMap[typeutil.String(key)] = value
			}
			return nil
		}
	}

	switch dataM.Kind() {
	case reflect.Map:
		if isDelete {
			dataM.SetMapIndex(
				reflect.ValueOf(key),
				reflect.Value{},
			)
		} else {
			dataM.SetMapIndex(
				reflect.ValueOf(key),
				reflect.ValueOf(value),
			)
		}
	case reflect.Slice, reflect.Array:
		if isDelete {
			return fmt.Errorf("Array item deletion not implemented")
		} else if typeutil.IsInteger(key) {
			dataM.Index(int(typeutil.Int(key)))
		} else {
			return fmt.Errorf("cannot set non-integer array index %q", key)
		}
	}

	return nil
}

func DeepSet(data interface{}, path []string, value interface{}) interface{} {
	if len(path) == 0 {
		return data
	}

	var first = path[0]
	var rest = make([]string, 0)

	if len(path) > 1 {
		rest = path[1:]
	}

	//  Leaf Nodes: this is where the value we're setting actually gets set/appended
	if len(rest) == 0 {
		//  parent element is an array; set the correct index or append if the index is out of bounds
		if typeutil.IsArray(data) {
			dataArray := sliceutil.Sliceify(data)

			if curIndex := int(typeutil.Int(first)); typeutil.IsInteger(first) {
				if curIndex >= len(dataArray) {
					for add := len(dataArray); add <= curIndex; add++ {
						dataArray = append(dataArray, nil)
					}
				}

				if curIndex < len(dataArray) {
					dataArray[curIndex] = value
					return dataArray
				}
			}

		} else if typeutil.IsMap(data) {
			if err := Set(data, first, value); err == nil {
				return data
			}

		} else if typeutil.IsStruct(data) {
			// we only accept a pointer to a struct here
			if dV := reflect.ValueOf(data); dV.Kind() == reflect.Ptr {
				// make sure dV is the underlying struct Value
				if dE := dV.Elem(); dE.Kind() == reflect.Struct {
					dV = dE
				} else {
					return data
				}

				var dT = dV.Type()

				for i := 0; i < dT.NumField(); i++ {
					if fT := dT.Field(i); fT.Name == first {
						if fV := dV.Field(i); fV.IsValid() && fV.CanSet() {
							typeutil.SetValue(dV.Field(i), value)
						}

						break
					}
				}
			}

			return data
		}

	} else {
		//  Array Embedding: this is where non-terminal array-index key components are processed
		if typeutil.IsInteger(rest[0]) {
			if typeutil.IsMap(data) {
				//  is the value at `first' in the map isn't present or isn't an array, create it
				var curVal = Get(data, first)

				if typeutil.IsArray(curVal) {
					curVal = sliceutil.Sliceify(curVal)
				} else {
					curVal = make([]interface{}, 0)
					Set(data, first, curVal)
				}

				// recurse into our cool array and do awesome stuff with it
				if err := Set(data, first, DeepSet(curVal, rest, value)); err == nil {
					return data
				}
			}

		} else {
			//  Intermediate Map Processing
			//    this is where branch nodes get created and populated via recursion
			//    depending on the data type of the input `data', non-existent maps
			//    will be created and either set to `data[first]' (the map)
			//    or appended to `data[first]' (the array)
			if typeutil.IsArray(data) {
				var dataArray = sliceutil.Sliceify(data)

				if curIndex := int(typeutil.Int(first)); typeutil.IsInteger(first) {
					if curIndex >= len(dataArray) {
						for add := len(dataArray); add <= curIndex; add++ {
							dataArray = append(dataArray, make(map[string]interface{}))
						}
					}

					if curIndex < len(dataArray) {
						dataArray[curIndex] = DeepSet(dataArray[curIndex], rest, value)
						return dataArray
					}
				}

			} else if dataMap, ok := data.(map[string]interface{}); ok {
				//  handle good old fashioned maps-of-maps
				//  is the value at 'first' in the map isn't present or isn't a map, create it
				var curVal, _ = dataMap[first]

				if !typeutil.IsMap(curVal) {
					dataMap[first] = make(map[string]interface{})
					curVal, _ = dataMap[first]
				}

				dataMap[first] = DeepSet(dataMap[first], rest, value)
				return dataMap
			}
		}
	}

	return data
}

func Append(maps ...map[string]interface{}) map[string]interface{} {
	var out = make(map[string]interface{})

	for _, mapV := range maps {
		for k, v := range mapV {
			out[k] = v
		}
	}

	return out
}

func Pluck(sliceOfMaps interface{}, key []string) []interface{} {
	var rv = make([]interface{}, 0)

	if sliceOfMaps == nil {
		return rv
	}

	WalkStruct(sliceOfMaps, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf && len(path) > 1 {
			var shouldInclude bool

			for i, _ := range path {
				if i == 0 {
					continue
				} else if (i - 1) < len(key) {
					if key[i-1] == `*` || path[i] == key[i-1] {
						shouldInclude = true
						continue
					} else {
						shouldInclude = false
						break
					}
				}
			}

			if shouldInclude {
				rv = append(rv, value)
			}
		}

		return nil
	})

	return rv
}

// Recursively walk through the given map, calling walkFn for each intermediate and leaf value.
func Walk(input interface{}, walkFn WalkFunc) error {
	return walkGeneric(input, nil, walkFn, false, nil)
}

// Recursively walk through the given map, calling walkFn for each intermediate and leaf value.
// This form behaves identically to Walk(), except that it will also recurse into structs, calling
// walkFn for all intermediate structs and fields.
func WalkStruct(input interface{}, walkFn WalkFunc) error {
	return walkGeneric(input, nil, walkFn, true, nil)
}

func walkGeneric(parent interface{}, path []string, walkFn WalkFunc, includeStruct bool, seen []uintptr) error {
	if parent == nil {
		return nil
	}

	var parentV = reflect.ValueOf(parent)

	switch parentV.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		var pp = parentV.Pointer()

		for _, s := range seen {
			if s == pp {
				return nil
			}
		}

		seen = append(seen, pp)
	}

	if parentV.Kind() == reflect.Ptr {
		parentV = parentV.Elem()

		if parentV.IsValid() {
			parent = parentV.Interface()
		} else {
			return nil
		}
	}

	switch parentV.Kind() {
	case reflect.Map:
		if err := walkFn(parent, path, false); err != nil {
			return returnSkipOrErr(err)
		}

		for _, key := range parentV.MapKeys() {
			valueV := parentV.MapIndex(key)
			subpath := append(path, fmt.Sprintf("%v", key.Interface()))

			if err := walkGeneric(valueV.Interface(), subpath, walkFn, includeStruct, seen); err != nil {
				return returnSkipOrErr(err)
			}
		}

	case reflect.Slice, reflect.Array:
		if err := walkFn(parent, path, false); err != nil {
			return returnSkipOrErr(err)
		}

		for i := 0; i < parentV.Len(); i++ {
			valueV := parentV.Index(i)
			subpath := append(path, fmt.Sprintf("%v", i))

			if err := walkGeneric(valueV.Interface(), subpath, walkFn, includeStruct, seen); err != nil {
				return returnSkipOrErr(err)
			}
		}

	case reflect.Struct:
		if includeStruct {
			if err := walkFn(parent, path, false); err != nil {
				return returnSkipOrErr(err)
			}

			for i := 0; i < parentV.NumField(); i++ {
				fieldV := parentV.Type().Field(i)

				// only operate on exported fields
				if fieldV.PkgPath == `` {
					valueV := parentV.Field(i)
					var subpath []string

					// if this field is embedded, don't add it to the path list because
					// it should be considered a first-class member of the parent struct
					if fieldV.Anonymous {
						subpath = path
					} else if name := fieldNameFromReflect(fieldV); name != `-` {
						subpath = append(path, name)
					} else {
						continue
					}

					if err := walkGeneric(valueV.Interface(), subpath, walkFn, includeStruct, seen); err != nil {
						return returnSkipOrErr(err)
					}
				}
			}
		} else {
			return walkFn(parent, path, true)
		}

	default:
		return walkFn(parent, path, true)
	}

	return nil
}

// Recursively remove all zero and empty values from the given map.
func Compact(input map[string]interface{}) (map[string]interface{}, error) {
	var output = make(map[string]interface{})

	if err := Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		if !typeutil.IsEmpty(value) {
			if typeutil.IsArray(value) {
				DeepSet(output, path, value)
				return SkipDescendants
			} else if isLeaf {
				DeepSet(output, path, value)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return output, nil
}

// Recursively merge the contents of the second map into the first one and return the result.
func Merge(first interface{}, second interface{}, options ...MergeOption) (map[string]interface{}, error) {
	if first != nil && !typeutil.IsKind(first, reflect.Map) {
		return nil, fmt.Errorf("first argument must be a map, got %T", first)
	}

	if second != nil && !typeutil.IsKind(second, reflect.Map) {
		return nil, fmt.Errorf("second argument must be a map, got %T", second)
	}

	var output = make(map[string]interface{})

	if err := Walk(first, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			DeepSet(output, path, value)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := Walk(second, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			if value != nil {
				if currentValue := DeepGet(output, path, nil); currentValue == nil {
					DeepSet(output, path, value)
				} else {
					currentV := reflect.ValueOf(currentValue)

					switch currentV.Type().Kind() {
					case reflect.Slice, reflect.Array:
						newPath := append(path, fmt.Sprintf("%d", currentV.Len()))
						DeepSet(output, newPath, value)

					default:
						if MergeOptions(options).Has(AppendValues) {
							DeepSet(output, path, []interface{}{currentValue, value})
						} else {
							DeepSet(output, path, value)
						}
					}
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return output, nil
}

// Take the input map and convert all values to strings.
func Stringify(input map[string]interface{}) map[string]string {
	var output = make(map[string]string)

	for k, v := range input {
		if str, err := stringutil.ToString(v); err == nil {
			output[k] = str
		} else {
			output[k] = fmt.Sprintf("!#ERR<%v>", err)
		}
	}

	return output
}

// Recursively walk the given map, performing automatic type conversion on all leaf nodes.
func Autotype(input interface{}) map[string]interface{} {
	var output = make(map[string]interface{})

	if err := Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		// if we encounter a Variant, use its autotyping instead of walking the struct
		if valueVar, ok := value.(typeutil.Variant); ok {
			value = valueVar.Auto()

			if !typeutil.IsEmpty(value) {
				DeepSet(output, path, value)
			}

			return SkipDescendants
		} else if isLeaf {
			if !typeutil.IsEmpty(value) {
				DeepSet(output, path, stringutil.Autotype(value))
			}
		}

		return nil
	}); err != nil {
		panic(err.Error())
	}

	return output
}

// Performs a JSONPath query against the given object and returns the results.
// JSONPath description, syntax, and examples are available at http://goessner.net/articles/JsonPath/.
func JSONPath(data interface{}, query string) (interface{}, error) {
	return utils.JSONPath(data, query, true)
}

func apply(includeStruct bool, input interface{}, fn ApplyFunc) map[string]interface{} {
	var output = make(map[string]interface{})

	wfn := func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			if fn != nil {
				if out, ok := fn(path, value); ok {
					DeepSet(output, path, out)
					return nil
				}
			}

			DeepSet(output, path, value)
		}

		return nil
	}

	var err error

	if includeStruct {
		err = WalkStruct(input, wfn)
	} else {
		err = Walk(input, wfn)
	}

	if err != nil {
		panic(err.Error())
	}

	return output
}

// Recursively walk the given map, calling the ApplyFunc for each leaf value.  If the second
// return value from the function is true, that value in the struct will be replaced with the first
// return value.  If false, the value will be left as-is.
func Apply(input interface{}, fn ApplyFunc) map[string]interface{} {
	return apply(false, input, fn)
}

// The same as Apply(), but will descend into structs.
func ApplyStruct(input interface{}, fn ApplyFunc) map[string]interface{} {
	return apply(true, input, fn)
}

// Perform a deep copy of the given map.
func DeepCopy(input interface{}) map[string]interface{} {
	return apply(false, input, nil)
}

// Perform a deep copy of the given map or struct, returning a map.
func DeepCopyStruct(input interface{}) map[string]interface{} {
	return apply(true, input, nil)
}

// Convert the given value to a slice using typeutil.Slice, then return each element as a Map.
func SliceOfMaps(input interface{}) (maps []*Map) {
	for _, v := range typeutil.Slice(input) {
		maps = append(maps, M(v))
	}

	return
}

func returnSkipOrErr(err error) error {
	if err == SkipDescendants {
		return nil
	} else {
		return err
	}
}

func fieldNameFromReflect(field reflect.StructField) string {
	if tag := field.Tag.Get(UnmarshalStructTag); tag != `` {
		if name, _ := stringutil.SplitPair(tag, `,`); name != `` {
			return name
		}
	}

	return field.Name
}

// Format the given string in the same manner as fmt.Sprintf, except data items that are
// maps or Map objects will be expanded using special patterns in the format string. Deeply-nested
// map values can be referenced using a format string "${path.to.value}".  Missing keys will return
// an empty string, or a fallback value may be provided like so: "${path.to.value|fallback}".
// The value may also specify a standard fmt.Sprintf pattern with "${path.to.value:%02d}" (or
// "${path.to.value|fallback:%02d}" for fallback values.)  Finally, a special case for time.Time values
// allows for the format string to be passed to time.Format: "${path.to.time:%January 2, 2006 (3:04pm)}".
func Sprintf(format string, data ...interface{}) string {
	var params []interface{}

MatchLoop:
	for {
		m := rxutil.Match(rxMapFmt, format)

		if m == nil {
			break
		}

		caps := m.NamedCaptures()
		placeholder := caps[`fmt`]

		if placeholder == `` {
			placeholder = `%v`
		}

		for _, d := range data {
			dm := M(d)

			if tm, ok := dm.Get(caps[`key`]).Value.(time.Time); ok {
				tmfmt := strings.TrimPrefix(placeholder, `%`)

				if tmfmt == `v` {
					tmfmt = time.RFC3339
				}

				params = append(params, tm.Format(tmfmt))
				format = m.ReplaceGroup(1, `%s`)
				continue MatchLoop
			} else if v := dm.String(caps[`key`]); v != `` {
				params = append(params, v)
				format = m.ReplaceGroup(1, placeholder)
				continue MatchLoop
			}
		}

		params = append(params, caps[`fallback`])
		format = m.ReplaceGroup(1, placeholder)
	}

	return fmt.Sprintf(format, params...)
}

// Same as Sprintf, but prints its output to standard output.
func Printf(format string, data ...interface{}) {
	fmt.Print(Sprintf(format, data...))
}

// Same as Sprintf, but writes output to the given writer.
func Fprintf(w io.Writer, format string, data ...interface{}) {
	fmt.Fprint(w, Sprintf(format, data...))
}
