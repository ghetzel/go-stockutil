// Helpers for type inflection and simplifying working with Golang generic interface types
package typeutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghetzel/go-stockutil/utils"
)

type TypeConvertFunc = utils.TypeConvertFunc

var scs = spew.ConfigState{
	Indent:            `    `,
	DisableCapacities: true,
	SortKeys:          true,
}

// Register's a handler used for converting one type to another. Type are checked in the following
// manner:  The input value's reflect.Type String() value is matched, falling back to its
// reflect.Kind String() value, finally checking for a special "*" value that matches any type.
// If the handler function returns nil, its value replaces the input value.  If the special error
// type PassthroughType is returned, the original value is returned unmodified.
func RegisterTypeHandler(handler TypeConvertFunc, types ...string) {
	utils.RegisterTypeHandler(handler, types...)
}

// Returns whether the given value represents the underlying type's zero value
func IsZero(value any) bool {
	return utils.IsZero(value)
}

// Returns whether the given value is "empty" in the semantic sense. Zero values
// are considered empty, as are arrays, slices, and maps containing only empty
// values (called recursively). Finally, strings are trimmed of whitespace and
// considered empty if the result is zero-length.
func IsEmpty(value any) bool {
	var valueV = reflect.ValueOf(value)

	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}

	// short circuit for zero values of certain types
	switch valueV.Kind() {
	case reflect.Struct:
		if IsZero(value) {
			return true
		}
	}

	switch valueV.Kind() {
	case reflect.Array, reflect.Slice:
		if valueV.Len() == 0 {
			return true
		} else {
			for i := 0; i < valueV.Len(); i++ {
				if indexV := valueV.Index(i); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
					return false
				}
			}

			return true
		}

	case reflect.Map:
		if valueV.Len() == 0 {
			return true
		} else {
			for _, keyV := range valueV.MapKeys() {
				if indexV := valueV.MapIndex(keyV); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
					return false
				}
			}

			return true
		}

	case reflect.Chan:
		if valueV.Len() == 0 {
			return true
		}

	case reflect.String:
		if len(strings.TrimSpace(fmt.Sprintf("%v", value))) == 0 {
			return true
		}
	}

	return false
}

// Return the concrete value pointed to by a pointer type, or within an
// interface type.  Allows functions receiving pointers to supported types
// to work with those types without doing reflection.
func ResolveValue(in any) any {
	if inV, ok := in.(Variant); ok {
		in = inV.Value
	}

	return utils.ResolveValue(in)
}

// Dectect whether the concrete underlying value of the given input is one or more
// Kinds of value.
func IsKind(in any, kinds ...reflect.Kind) bool {
	return utils.IsKind(in, kinds...)
}

// Return whether the given input is a discrete scalar value (ints, floats, bools,
// strings), otherwise known as "primitive types" in some other languages.
func IsScalar(in any) bool {
	if !IsKind(in, utils.CompoundTypes...) {
		return true
	}

	return false
}

// Returns whether the given value is a slice or array.
func IsArray(in any) bool {
	return IsKind(in, utils.SliceTypes...)
}

// Returns whether the given value is a map.
func IsMap(in any) bool {
	return IsKind(in, reflect.Map)
}

// Returns whether the given value is a struct.
func IsStruct(in any) bool {
	return IsKind(in, reflect.Struct)
}

// Returns whether the given value is a function of any kind
func IsFunction(in any) bool {
	return IsKind(in, reflect.Func)
}

// Returns whether the given value represents a numeric value.
func IsNumeric(in any) bool {
	return utils.IsNumeric(in)
}

// Returns whether the given value represents an integer value.
func IsInteger(in any) bool {
	return utils.IsInteger(in)
}

// Returns whether the given value represents a floating point value.
func IsFloat(in any) bool {
	return utils.IsFloat(in)
}

// Return whether the value can be interpreted as a time.
func IsTime(in any) bool {
	return VV(in).IsTime()
}

// Return whether the value can be interpreted as a duration.
func IsDuration(in any) bool {
	return VV(in).IsDuration()
}

// Returns whether the given value is a function.  If inParams is not -1, the function must
// accept that number of arguments.  If outParams is not -1, the function must return that
// number of values.
func IsFunctionArity(in any, inParams int, outParams int) bool {
	if IsKind(in, reflect.Func) {
		var inT = reflect.TypeOf(in)

		if inParams < 0 || inParams >= 0 && inT.NumIn() == inParams {
			if outParams < 0 || outParams >= 0 && inT.NumOut() == outParams {
				return true
			}
		}
	}

	return false
}

// Returns the length of the given value that could have a length (strings, slices, arrays,
// maps, and channels).  If the value is not a type that has a length, -1 is returned.
func Len(in any) int {
	if IsKind(in, reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String) {
		return reflect.ValueOf(in).Len()
	} else {
		return -1
	}
}

// Returns a pretty-printed string representation of the given values.
func Dump(in1 any, in ...any) string {
	return scs.Sdump(append([]any{in1}, in...)...)
}

// Returns a pretty-printed string representation of the given values.
func Dumpf(format string, in ...any) string {
	return fmt.Sprintf(format, scs.Sdump(in...))
}

// Attempts to set the given reflect.Value to the given interface value
func SetValue(target any, value any) error {
	var targetV, valueV, originalV reflect.Value

	// if we were given a reflect.Value target, then we shouldn't take the reflect.ValueOf that
	if tV, ok := target.(reflect.Value); ok {
		targetV = tV
	} else {
		targetV = reflect.ValueOf(target)

		if targetV.Kind() == reflect.Struct {
			return fmt.Errorf("Must pass a pointer to a struct instance, got %T", target)
		} else if targetV.Kind() == reflect.Ptr {
			// dereference pointers to get at the real destination
			targetV = targetV.Elem()
		}
	}

	// targets must be valid in order to set them to values
	if !targetV.IsValid() {
		return fmt.Errorf("Target %T is not valid", target)
	}

	// perform custom type conversions (if any)
	if v, err := utils.ConvertCustomType(value); err == nil {
		value = v
	} else if err != utils.PassthroughType {
		return err
	}

	// if the value we were given was a reflect.Value, just use that
	if vV, ok := value.(reflect.Value); ok {
		originalV = vV
		valueV = vV
	} else {
		originalV = reflect.ValueOf(value)
		valueV = reflect.ValueOf(ResolveValue(value))
	}

	// get the underlying value that was passed in
	if valueV.IsValid() {
		var targetT = targetV.Type()
		var valueT = valueV.Type()

		// if the target value is a string-a-like, stringify whatever we got in
		if targetT.Kind() == reflect.String && valueV.CanInterface() {
			valueV = reflect.ValueOf(fmt.Sprintf("%v", valueV.Interface()))
			valueT = valueV.Type()

			if !valueV.IsValid() {
				return fmt.Errorf(
					"Converting %T to %v produced an invalid value",
					value,
					targetT,
				)
			}
		}

		if valueT.AssignableTo(targetT) {
			// attempt direct assignment
			targetV.Set(valueV)
		} else if valueT.ConvertibleTo(targetT) {
			// attempt type conversion
			targetV.Set(valueV.Convert(targetT))
		} else if targetV.Kind() == reflect.Ptr {
			if originalV.Kind() == reflect.Ptr {
				return SetValue(targetV, originalV)
			} else {
				return fmt.Errorf(
					"Unable to set target: value for target %v must be given as a pointer",
					targetT,
				)
			}
		} else {
			switch kind := targetV.Kind(); kind {
			case reflect.Struct:
				if embeddedV := targetV.FieldByName(valueT.Name()); embeddedV.IsValid() {
					if err := SetValue(embeddedV, value); err == nil {
						return nil
					}
				}

			case reflect.Array, reflect.Slice:
				if IsArray(value) {
					var valueA = utils.Sliceify(value)
					var repl = targetV

					if targetV.Len() < len(valueA) {
						repl = reflect.MakeSlice(targetT, len(valueA), len(valueA))
					}

					for i := 0; i < repl.Len(); i++ {
						var elem = repl.Index(i)

						if err := SetValue(elem, valueA[i]); err != nil {
							return fmt.Errorf("cannot set index %d: %v", i, err)
						}
					}

					return SetValue(targetV, repl)
				}
			}

			// handle some well-know, type specific edge cases
			switch value.(type) {
			case time.Time:
				return SetValue(target, value.(time.Time).UnixNano())
			case *time.Time:
				if tm := value.(*time.Time); tm != nil {
					return SetValue(target, tm.UnixNano())
				}
			}

			// no dice.
			return fmt.Errorf(
				"Unable to set target: %T has no path to becoming %v",
				value,
				targetT,
			)
		}
	} else {
		return fmt.Errorf("Unable to set target to the given %T value", value)
	}

	return nil
}

func IsKindOfInteger(in any) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(
		kind,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
	)
}

func IsKindOfFloat(in any) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(
		kind,
		reflect.Float32,
		reflect.Float64,
	)
}

func IsKindOfBool(in any) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(kind, reflect.Bool)
}

func IsKindOfString(in any) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else if inT := reflect.TypeOf(in); inT != nil {
		kind = inT.Kind()
	}

	return IsKind(kind, reflect.String)
}

// Provide a variable to encode as JSON, and an optional indent string.  If no indent argument is
// provided, the default indent is "  " (two spaces).  If an empty string is explcitly provided
// for the indent argument, the output will not be indented (single line).
func JSON(in any, indent ...string) string {
	var i string
	var out []byte
	var err error

	if len(indent) > 0 {
		i = indent[0]
	} else {
		i = `  `
	}

	if i == `` {
		out, err = json.Marshal(in)
	} else {
		out, err = json.MarshalIndent(in, ``, i)
	}

	if err == nil {
		return string(out)
	} else {
		return ``
	}
}
