package utils

import (
	"fmt"
	"reflect"
)

var Stop = fmt.Errorf("stop iterating")
var EachChanMaxItems = 1048576

type IterationFunc func(i int, value any) error

// Iterate through each element of the given array, slice or channel; calling
// iterFn exactly once for each element.  Otherwise, call iterFn one time
// with the given input as the argument.
func SliceEach(slice any, iterFn IterationFunc, preserve ...reflect.Kind) error {
	if iterFn == nil {
		return nil
	}

	slice = ResolveValue(slice)

	if inT := reflect.TypeOf(slice); inT != nil {
		switch inT.Kind() {
		case reflect.Slice, reflect.Array:
			var sliceV = reflect.ValueOf(slice)

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

			var mapV = reflect.ValueOf(slice)

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

			var structV = reflect.ValueOf(slice)

			for i := 0; i < structV.Type().NumField(); i++ {
				var field = structV.Type().Field(i)

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

		case reflect.Chan:
			var sliceC = reflect.ValueOf(slice)
			var i int

			for {
				if item, ok := sliceC.Recv(); ok {
					if item.IsValid() && item.CanInterface() {
						if err := iterFn(i, item.Interface()); err != nil {
							if err == Stop {
								return nil
							} else {
								return err
							}
						} else {
							i++
						}
					}
				} else if i >= EachChanMaxItems {
					break
				} else {
					break
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

// Takes some input value and returns it as a slice.
func Sliceify(in any) []any {
	if in == nil {
		return nil
	}

	var out = make([]any, 0)

	SliceEach(in, func(_ int, v any) error {
		out = append(out, v)
		return nil
	})

	return out
}
