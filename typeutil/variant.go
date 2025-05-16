package typeutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/ghetzel/go-stockutil/utils"
)

// Represents an interface type with helper functions for making it easy to do
// type conversions.
type Variant struct {
	Value any
}

// Shortcut for creating a Variant.
func V(value any) Variant {
	if v, ok := value.(Variant); ok {
		return v
	} else {
		return Variant{
			Value: value,
		}
	}
}

// Returns a pointer to a variant.
func VV(value any) *Variant {
	var v = V(value)
	return &v
}

// Returns whether the underlying value is nil.
func (self Variant) IsNil() bool {
	return (self.Value == nil)
}

// Returns whether the underlying value is a zero value.
func (self Variant) IsZero() bool {
	return IsZero(self.Value)
}

// Return whether the value is of the given reflect.Kind
func (self *Variant) IsKind(kind reflect.Kind) bool {
	if vV, ok := self.Value.(reflect.Value); ok && vV.IsValid() {
		return (vV.Type().Kind() == kind)
	} else if vT, ok := self.Value.(reflect.Type); ok && vT != nil {
		return (vT.Kind() == kind)
	} else if vT := reflect.TypeOf(self.Value); vT != nil {
		return (vT.Kind() == kind)
	} else {
		return false
	}
}

// Return the value as a string, or an empty string if the value could not be converted.
func (self Variant) String() string {
	if v, err := utils.ConvertToString(self.Value); err == nil {
		return v
	} else {
		return ``
	}
}

// Return true if the value can be interpreted as a boolean true value, or false otherwise.
func (self Variant) Bool() bool {
	if v, err := utils.ConvertToBool(self.Value); err == nil {
		return v
	} else if vF, err := utils.ConvertToFloat(self.Value); err == nil && vF == 0 {
		return false
	} else {
		// use a more relaxed set of values for determining "true" because
		// the user has very explicitly asked us to try
		truthy, _ := template.IsTrue(self.Value)
		return truthy
	}
}

// Return the value as a float if it can be interpreted as such, or 0 otherwise.
func (self Variant) Float() float64 {
	if v, err := utils.ConvertToFloat(self.Value); err == nil {
		return v
	} else {
		return 0
	}
}

// Return the value as an integer if it can be interpreted as such, or 0 otherwise. Float values
// will be truncated to integers.
func (self Variant) Int() int64 {
	if v, err := utils.ConvertToFloat(self.Value); err == nil {
		return int64(v)
	} else {
		return 0
	}
}

// Return the value as a native integer if it can be interpreted as such, or 0 otherwise.
// Float values will be truncated to integers.
func (self Variant) NInt() int {
	return int(self.Int())
}

// Return the value as a slice of Variants. Scalar types will return a slice containing
// a single Variant element representing the value.
func (self Variant) Slice() []Variant {
	var values = make([]Variant, 0)

	utils.SliceEach(self.Value, func(_ int, v any) error {
		values = append(values, Variant{
			Value: v,
		})
		return nil
	})

	return values
}

// Same as Slice(), but returns a []string.
func (self Variant) Strings() []string {
	var values = self.Slice()
	var out = make([]string, len(values))

	for i, value := range values {
		out[i] = value.String()
	}

	return out
}

// Converts the value to a string, then splits on the given delimiter.
func (self Variant) Split(on string) []string {
	if s := self.String(); s == `` {
		return nil
	} else {
		return strings.Split(s, on)
	}
}

// Return the value converted to an error, or nil if it is not an error.
func (self Variant) Err() error {
	if err, ok := self.Value.(error); ok {
		return err
	} else {
		return nil
	}
}

// Return the value automatically converted to the appropriate type.
func (self Variant) Auto() any {
	return utils.Autotype(self.Value)
}

// Return the value as a time.Time if it can be interpreted as such, or zero time otherwise.
func (self Variant) Time() time.Time {
	if v, err := utils.ConvertToTime(self.Value); err == nil {
		return v
	} else {
		return time.Time{}
	}
}

// Return the value as a time.Duration if it can be interpreted as such, or zero otherwise.
func (self Variant) Duration() time.Duration {
	if v, err := utils.ParseDuration(self.String()); err == nil {
		return v
	} else {
		return 0
	}
}

// Return the value at key as a byte slice.
func (self Variant) Bytes() []byte {
	if v, err := utils.ConvertToBytes(self.Value); err == nil {
		return v
	} else {
		return []byte{}
	}
}

func (self Variant) Interface() any {
	var out = self.Value

	for {
		if subvariant, ok := out.(Variant); ok {
			out = subvariant.Value
		} else if subvariant, ok := out.(*Variant); ok {
			out = subvariant.Value
		} else {
			break
		}
	}

	return out
}

// Return the value as a map[Variant]Variant if it can be interpreted as such, or an empty map otherwise.
// If the variant contains a struct, and a tagName is specified, the key names of the output map will be taken
// from the struct field's tag value, consistent with the rules used in encoding/json.
func (self Variant) Map(tagName ...string) map[Variant]Variant {
	var output = make(map[Variant]Variant)
	var tn string

	if len(tagName) > 0 && tagName[0] != `` {
		tn = tagName[0]
	}

	if IsMap(self.Value) {
		var mapV = reflect.ValueOf(self.Value)

		for _, key := range mapV.MapKeys() {
			if key.CanInterface() {
				if value := mapV.MapIndex(key); value.CanInterface() {
					output[V(key.Interface())] = V(value.Interface())
				}
			}
		}
	} else if elem := ResolveValue(self.Value); IsStruct(elem) {
		var structV = reflect.ValueOf(elem)
		var structT = structV.Type()

	FieldLoop:
		for i := 0; i < structT.NumField(); i++ {
			if structF := structT.Field(i); !structF.Anonymous {
				if value := structV.Field(i); value.IsValid() {
					var key = structF.Name

					if tn != `` {
						var parts = strings.Split(structF.Tag.Get(tn), `,`)

						if tag := parts[0]; tag != `` {
							key = tag
						}

						if len(parts) > 1 {
							for _, p := range parts {
								switch p {
								case `omitempty`:
									if IsZero(value) {
										continue FieldLoop
									}
								}
							}
						}
					}

					if value.CanInterface() {
						output[V(key)] = V(value.Interface())
					}
				}
			}
		}
	}

	return output
}

// Return the value as a map[string]any if it can be interpreted as such, or an empty map otherwise.
func (self Variant) MapNative(tagName ...string) map[string]any {
	var out = make(map[string]any)

	for k, v := range self.Map(tagName...) {
		out[k.String()] = v.Interface()
	}

	return out
}

// Satisfy the json.Marshaler interface
func (self Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Auto())
}

// Return whether the value is an array/slice type.
func (self *Variant) IsArray() bool {
	return IsArray(self.Value)
}

// Return whether the value is a map type.
func (self *Variant) IsMap() bool {
	return IsMap(self.Value)
}

// Return whether the value is a scalar type.
func (self *Variant) IsScalar() bool {
	return IsScalar(self.Value)
}

// Return whether the value can be interpreted as a real number.
func (self *Variant) IsNumeric() bool {
	return IsNumeric(self.Value)
}

// Return whether the value can be interpreted as a time.
func (self *Variant) IsTime() bool {
	return !self.Time().IsZero()
}

// Return whether the value can be interpreted as a duration.
func (self *Variant) IsDuration() bool {
	return (self.Duration() != 0)
}

func (self *Variant) Append(values ...any) error {
	var base []any

	if self.IsArray() {
		base = utils.Sliceify(self.Value)
	} else {
		base = []any{self.Value}
	}

	for _, val := range values {
		for _, subval := range utils.Sliceify(val) {
			if valV, ok := subval.(reflect.Value); ok {
				if valV.IsValid() && valV.CanInterface() {
					base = append(base, valV.Interface())
				} else {
					return fmt.Errorf("invalid value")
				}
			} else {
				base = append(base, subval)
			}
		}
	}

	self.Value = base
	return nil
}

func (self Variant) OrString(or ...any) string {
	if v := self.String(); v != `` {
		return v
	}

	for _, orval := range or {
		if v := String(orval); v != `` {
			return v
		}
	}

	return ``
}

func (self Variant) OrBool(or ...any) bool {
	if self.Value == nil {
		for _, orval := range or {
			if Bool(orval) {
				return true
			}
		}

		return false
	} else {
		return self.Bool()
	}
}

func (self Variant) OrFloat(or ...any) float64 {
	if v := self.Float(); v != 0 {
		return v
	}

	for _, orval := range or {
		if v := Float(orval); v != 0 {
			return v
		}
	}

	return 0
}

func (self Variant) OrInt(or ...any) int64 {
	if v := self.Int(); v != 0 {
		return v
	}

	for _, orval := range or {
		if v := Int(orval); v != 0 {
			return v
		}
	}

	return 0
}

func (self Variant) OrNInt(or ...any) int {
	if v := self.NInt(); v != 0 {
		return v
	}

	for _, orval := range or {
		if v := NInt(orval); v != 0 {
			return v
		}
	}

	return 0
}

func (self Variant) OrAuto(or ...any) any {
	if v := self.Auto(); !IsZero(v) {
		return v
	}

	for _, orval := range or {
		if !IsZero(orval) {
			return orval
		}
	}

	return nil
}

func (self Variant) OrTime(or ...any) time.Time {
	if v := self.Time(); !v.IsZero() {
		return v
	}

	for _, orval := range or {
		if ov := V(orval).Time(); !ov.IsZero() {
			return ov
		}
	}

	return time.Time{}
}

func (self Variant) OrDuration(or ...any) time.Duration {
	if v := self.Duration(); v != 0 {
		return v
	}

	for _, orval := range or {
		if ov := V(orval).Duration(); ov != 0 {
			return ov
		}
	}

	return 0
}

func (self Variant) OrBytes(or ...[]byte) []byte {
	if v := self.Bytes(); len(v) > 0 {
		return v
	}

	for _, orval := range or {
		if len(orval) > 0 {
			return orval
		}
	}

	return nil
}

func (self Variant) Or(or ...any) any {
	if v := self.Interface(); !IsZero(v) {
		return v
	}

	for _, orval := range or {
		if !IsZero(orval) {
			return orval
		}
	}

	return nil
}

// IsLessThan reports whether the given value should sort before the current variant value, taking special
// care to compare like types appropriately, such as detecting numbers and performing a numeric comparison,
// or detecting dates, times, and durations and comparing them temporally.
func (self Variant) IsLessThan(j any) bool {
	var jV Variant = V(j)

	// we're going to say that if we encounter a nil, we always want that nil to sort before this
	if jV.IsNil() {
		return false
	}

	if self.IsNumeric() && jV.IsNumeric() {
		// compare numeric values numerically
		return self.Float() < jV.Float()
	} else if self.IsTime() && jV.IsTime() {
		// compare times
		return self.Time().Before(jV.Time())
	} else if self.IsTime() && jV.IsTime() {
		// compare durations
		return self.Duration() < jV.Duration()
	} else {
		// fallback to lexical comparison
		return self.String() < jV.String()
	}
}

func (self Variant) IsFunction(signature ...string) bool {
	if !self.IsKind(reflect.Func) {
		return false
	}

	if len(signature) > 0 {
		return (FunctionMatchesSignature(self.Value, signature[0]) == nil)
	} else {
		return true
	}
}

// Package-level string converter
func String(in any) string {
	return V(in).String()
}

// Package-level bool converter
func Bool(in any) bool {
	return V(in).Bool()
}

// Package-level float converter
func Float(in any) float64 {
	return V(in).Float()
}

// Package-level int64 converter
func Int(in any) int64 {
	return V(in).Int()
}

// Package-level native int converter
func NInt(in any) int {
	return V(in).NInt()
}

// Package-level slice converter
func Slice(in any) []Variant {
	return V(in).Slice()
}

// Package-level string slice converter
func Strings(in any) []string {
	return V(in).Strings()
}

// Package-level string splitter.
func Split(in any, on string) []string {
	return V(in).Split(on)
}

// Package-level error converter
func Err(in any) error {
	return V(in).Err()
}

// Package-level auto converter
func Auto(in any) any {
	return V(in).Auto()
}

// Package-level time converter
func Time(in any) time.Time {
	return V(in).Time()
}

// Package-level duration converter
func Duration(in any) time.Duration {
	return V(in).Duration()
}

// Package-level bytes converter
func Bytes(in any) []byte {
	return V(in).Bytes()
}

// Package-level map converter
func Map(in any, tagName ...string) map[Variant]Variant {
	return V(in).Map(tagName...)
}

// Package-level map[string]any converter
func MapNative(in any, tagName ...string) map[string]any {
	return V(in).MapNative(tagName...)
}

// Return a new Variant with a nil value.
func Nil() Variant {
	return V(nil)
}

func OrString(first any, rest ...any) string {
	return V(first).OrString(rest...)
}

func OrBool(first any, rest ...any) bool {
	return V(first).OrBool(rest...)
}

func OrFloat(first any, rest ...any) float64 {
	return V(first).OrFloat(rest...)
}

func OrInt(first any, rest ...any) int64 {
	return V(first).OrInt(rest...)
}

func OrNInt(first any, rest ...any) int {
	return V(first).OrNInt(rest...)
}

func OrAuto(first any, rest ...any) any {
	return V(first).OrAuto(rest...)
}

func OrTime(first any, rest ...any) time.Time {
	return V(first).OrTime(rest...)
}

func OrDuration(first any, rest ...any) time.Duration {
	return V(first).OrDuration(rest...)
}

func OrBytes(first []byte, rest ...[]byte) []byte {
	return V(first).OrBytes(rest...)
}

func IsLessThan(a any, b any) bool {
	return VV(a).IsLessThan(b)
}
