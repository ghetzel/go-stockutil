// Utilities for working with and manipulating structs.
package structutil

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ghetzel/go-stockutil/typeutil"
)

var StopIterating = errors.New(`stop iterating`)

// Receives a struct field name, the value of that field in the source struct, and the value for that field in the destination struct.
// Returns the value that should be placed in the destination struct fields.  If the returned bool is false, no changes will
// be made.
type StructValueFunc func(field string, sourceValue any, destValue any) (any, bool)
type StructFieldFunc func(field *reflect.StructField, value reflect.Value) error

// Iterates over all exported and embedded fields in the given struct, calling fn for each field.
func FieldsFunc(in any, fn StructFieldFunc) error {
	if in == nil || fn == nil {
		return nil
	}

	var inValu reflect.Value

	if v, ok := in.(reflect.Value); ok {
		inValu = v
	} else {
		inValu = reflect.ValueOf(in)
	}

	if inValu.Kind() == reflect.Ptr {
		inValu = inValu.Elem()
	}

	var inType = inValu.Type()

FieldLoop:
	for i := 0; i < inType.NumField(); i++ {
		var fieldT = inType.Field(i)
		var fieldV = inValu.Field(i)

		// recursively include embedded fields
		if fieldT.Anonymous {
			if err := FieldsFunc(fieldV, fn); err != nil {
				return err
			}
		}

		// only exported field names leave this empty, so skip if it's not (i.e.: we have an unexported field)
		if fieldT.PkgPath != `` {
			continue
		}

		switch err := fn(&fieldT, fieldV); err {
		case StopIterating:
			break FieldLoop
		case nil:
			continue
		default:
			return err
		}
	}

	return nil
}

func CopyFunc(dest any, source any, fn StructValueFunc) error {
	if dest == nil || source == nil || fn == nil {
		return nil
	}

	var destV reflect.Value
	var srcV reflect.Value

	if dV, ok := dest.(reflect.Value); ok {
		destV = dV
	} else {
		destV = reflect.ValueOf(dest)
	}

	if sV, ok := source.(reflect.Value); ok {
		srcV = sV
	} else {
		srcV = reflect.ValueOf(source)
	}

	if dV, err := validatePtrToStruct(`destination`, destV); err == nil {
		destV = dV
	} else {
		return err
	}

	if sV, err := validatePtrToStruct(`source`, srcV); err == nil {
		srcV = sV
	} else {
		return err
	}

	var destT = destV.Type()
	var srcT = srcV.Type()

	for s := 0; s < srcT.NumField(); s++ {
		var sFieldT = srcT.Field(s)
		var sFieldV = srcV.Field(s)

		// only exported field names leave this empty, so skip if it's not (i.e.: we have an unexported field)
		if sFieldT.PkgPath != `` {
			continue
		}

		if dFieldT, ok := destT.FieldByName(sFieldT.Name); ok {
			var dFieldV = destV.FieldByName(dFieldT.Name)

			if dFieldT.Anonymous {
				if err := CopyFunc(dFieldV, sFieldV, fn); err != nil {
					return err
				}
			} else {
				if sFieldV.CanInterface() && dFieldV.CanInterface() {
					if repl, ok := fn(dFieldT.Name, sFieldV.Interface(), dFieldV.Interface()); ok {
						if dFieldV.CanSet() {
							if err := typeutil.SetValue(dFieldV, repl); err != nil {
								return err
							}
						} else {
							return fmt.Errorf("field %q is not settable", dFieldT.Name)
						}
					}
				} else {
					return fmt.Errorf("Cannot retrieve field value %q", dFieldT.Name)
				}
			}
		}
	}

	return nil
}

// Copy all values from the source into the destination, provided the source value for the corresponding
// field is not that type's zero value.
func CopyNonZero(dest any, source any) error {
	return CopyFunc(dest, source, func(name string, s any, d any) (any, bool) {
		if typeutil.IsZero(s) {
			return nil, false
		} else {
			return s, true
		}
	})
}

func validatePtrToStruct(name string, obj reflect.Value) (reflect.Value, error) {
	if obj.Kind() == reflect.Ptr {
		if obj.Elem().Kind() == reflect.Struct {
			return obj.Elem(), nil
		} else {
			return reflect.Value{}, fmt.Errorf("bad %s: expected pointer to struct, got %v", name, obj.Type())
		}
	} else {
		return reflect.Value{}, fmt.Errorf("bad %s: expected pointer to struct, got %v", name, obj.Type())
	}
}
