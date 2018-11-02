package utils

import (
	"fmt"
	"reflect"
)

var Stop = fmt.Errorf("stop iterating")

type IterationFunc func(i int, value interface{}) error

// Iterate through each element of the given array or slice, calling
// iterFn exactly once for each element.  Otherwise, call iterFn one time
// with the given input as the argument.
//
func SliceEach(slice interface{}, iterFn IterationFunc, preserve ...reflect.Kind) error {
	if iterFn == nil {
		return nil
	}

	slice = ResolveValue(slice)

	if inT := reflect.TypeOf(slice); inT != nil {
		switch inT.Kind() {
		case reflect.Slice, reflect.Array:
			sliceV := reflect.ValueOf(slice)

			for i := 0; i < sliceV.Len(); i++ {
				if err := iterFn(i, sliceV.Index(i).Interface()); err != nil {
					if err == Stop {
						return nil
					} else {
						return err
					}
				}
			}
		case reflect.Map:
			for _, p := range preserve {
				if p == reflect.Map {
					if err := iterFn(0, slice); err != nil {
						if err == Stop {
							return nil
						} else {
							return err
						}
					} else {
						return nil
					}
				}
			}

			mapV := reflect.ValueOf(slice)

			for i, key := range mapV.MapKeys() {
				if valueV := mapV.MapIndex(key); valueV.IsValid() && valueV.CanInterface() {
					if err := iterFn(i, valueV.Interface()); err != nil {
						if err == Stop {
							return nil
						} else {
							return err
						}
					}
				}
			}

		case reflect.Struct:
			for _, p := range preserve {
				if p == reflect.Struct {
					if err := iterFn(0, slice); err != nil {
						if err == Stop {
							return nil
						} else {
							return err
						}
					} else {
						return nil
					}
				}
			}

			structV := reflect.ValueOf(slice)

			for i := 0; i < structV.Type().NumField(); i++ {
				field := structV.Type().Field(i)

				if field.Name != `` {
					if valueV := structV.Field(i); valueV.IsValid() && valueV.CanInterface() {
						if err := iterFn(i, valueV.Interface()); err != nil {
							if err == Stop {
								return nil
							} else {
								return err
							}
						}
					}
				}
			}

		default:
			if err := iterFn(0, slice); err != nil {
				if err == Stop {
					return nil
				} else {
					return err
				}
			}
		}
	}

	return nil
}
