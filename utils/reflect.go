package utils

import (
	"reflect"
)

var SliceTypes = []reflect.Kind{
	reflect.Slice,
	reflect.Array,
}

var CompoundTypes = []reflect.Kind{
	reflect.Invalid,
	reflect.Complex64,
	reflect.Complex128,
	reflect.Array,
	reflect.Chan,
	reflect.Func,
	reflect.Interface,
	reflect.Map,
	reflect.Ptr,
	reflect.Slice,
	reflect.Struct,
}

func ResolveValue(in any) any {
	var inV reflect.Value

	if vV, ok := in.(reflect.Value); ok {
		inV = vV
	} else {
		inV = reflect.ValueOf(in)
	}

	if inV.IsValid() {
		if inT := inV.Type(); inT == nil {
			return nil
		}

		switch inV.Kind() {
		case reflect.Ptr, reflect.Interface:
			return ResolveValue(inV.Elem())
		}

		in = inV.Interface()
	}

	return in
}

// Dectect whether the concrete underlying value of the given input is one or more
// Kinds of value.
func IsKind(in any, kinds ...reflect.Kind) bool {
	var inT reflect.Type

	if v, ok := in.(reflect.Value); ok && v.IsValid() {
		inT = v.Type()
	} else if v, ok := in.(reflect.Type); ok {
		inT = v
	} else {
		in = ResolveValue(in)
		inT = reflect.TypeOf(in)
	}

	if inT == nil {
		return false
	}

	for _, k := range kinds {
		if inT.Kind() == k {
			return true
		}
	}

	return false
}

// Returns whether the given value represents the underlying type's zero value
func IsZero(value any) bool {
	if value == nil {
		return true
	} else if valueV, ok := value.(reflect.Value); ok && valueV.IsValid() {
		if valueV.CanInterface() {
			value = valueV.Interface()
		}
	}

	return reflect.DeepEqual(
		value,
		reflect.Zero(reflect.TypeOf(value)).Interface(),
	)
}
