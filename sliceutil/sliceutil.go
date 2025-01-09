// Utilities for converting, manipulating, and iterating over slices
package sliceutil

import (
	"reflect"
	"strings"

	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/go-stockutil/utils"
	"github.com/juliangruber/go-intersect"
)

var Stop = utils.Stop

type IterationFunc = utils.IterationFunc
type CompareFunc func(i int, first interface{}, second interface{}) bool // {}
type MapFunc func(i int, value interface{}) interface{}
type MapStringFunc func(i int, value string) string

var StrictEqualityCompare = func(_ int, first interface{}, second interface{}) bool {
	if first == second {
		return true
	}

	return false
}

var RelaxedEqualityCompare = func(_ int, first interface{}, second interface{}) bool {
	if v, err := stringutil.RelaxedEqual(first, second); err == nil && v == true {
		return true
	}

	return false
}

// Return whether the given slice contains the given value.  If a comparator is provided, it will
// be used to compare the elements.
func Contains(in interface{}, value interface{}, comparators ...CompareFunc) bool {
	if len(comparators) == 0 {
		comparators = []CompareFunc{StrictEqualityCompare}
	}

	comparator := comparators[0]

	if inV := reflect.ValueOf(in); inV.IsValid() {
		for i := 0; i < inV.Len(); i++ {
			if current := inV.Index(i); current.IsValid() {
				if comparator(i, value, current.Interface()) {
					return true
				}
			}
		}
	}

	return false
}

// Returns whether the given string slice contains a given string.
func ContainsString(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}

	return false
}

// Returns whether the given string slice contains any of the following strings.
func ContainsAnyString(list []string, elems ...string) bool {
	for _, e := range elems {
		if ContainsString(list, e) {
			return true
		}
	}

	return false
}

// Returns whether the given string slice contains all of the following strings.
func ContainsAllStrings(list []string, elems ...string) bool {
	for _, e := range elems {
		if !ContainsString(list, e) {
			return false
		}
	}

	return true
}

// Return the intersection of two string slices.
func IntersectStrings(a []string, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return make([]string, 0)
	}

	return Stringify(intersect.Simple(a, b))
}

// Return the intersection of two slices.
func Intersect(a interface{}, b interface{}) []interface{} {
	return Sliceify(intersect.Simple(a, b))
}

// Return the slice that results from removing elements in second from the first.
func Difference(first interface{}, second interface{}) []interface{} {
	var diff = make(map[interface{}]bool)
	var out = make([]interface{}, 0)
	var aS = Sliceify(first)
	var bS = Sliceify(second)

	if len(aS) == 0 {
		return out
	} else if len(bS) == 0 {
		return aS
	}

	for _, item := range bS {
		diff[item] = true
	}

	for _, item := range aS {
		if _, ok := diff[item]; !ok {
			out = append(out, item)
		}
	}

	return out
}

// Removes all elements from the given interface slice that are "empty", which is
// defined as being nil, a nil or zero-length array, chan, map, slice, or string.
//
// The zero values of any other type are not considered empty and will remain in
// the return value.
func Compact(in interface{}) []interface{} {
	if in == nil {
		return nil
	}

	rv := make([]interface{}, 0)

	for _, v := range Sliceify(in) {
		if v != nil {
			vV := reflect.ValueOf(v)

			switch vV.Kind() {
			case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
				if vV.Len() > 0 {
					rv = append(rv, v)
				}
			default:
				rv = append(rv, v)
			}
		}
	}

	return rv
}

// Returns the given slice as a single-level flattened array.
func Flatten(in interface{}) []interface{} {
	var out = make([]interface{}, 0)

	Each(in, func(_ int, value interface{}) error {
		if typeutil.IsArray(value) {
			out = append(out, Flatten(value)...)
		} else {
			out = append(out, value)
		}

		return nil
	})

	return out
}

// Removes all zero-length strings from the given string slice, returning a new
// slice with the values removed.
func CompactString(in []string) []string {
	if in == nil {
		return nil
	}

	rv := make([]string, 0)

	for _, v := range in {
		if v != `` {
			rv = append(rv, v)
		}
	}

	return rv
}

// Splits the given string by a delimiter, then eliminates any zero-length items.
func SplitCompact(in string, delimiter string) (out []string) {
	out = strings.Split(in, delimiter)
	out = CompactString(out)
	return
}

// Splits the given string by a delimiter, trims the space around each item,
// then eliminates any zero-length results.
func SplitTrimSpaceCompact(in string, delimiter string) (out []string) {
	out = strings.Split(in, delimiter)

	for i, v := range out {
		out[i] = strings.TrimSpace(v)
	}

	out = CompactString(out)
	
	return
}

// Converts all elements of the given interface slice to strings using the "%v"
// format string via the fmt package.
func Stringify(in interface{}) []string {
	if arr := Sliceify(in); len(arr) > 0 {
		out := make([]string, len(arr))

		for i, item := range arr {
			out[i] = typeutil.String(item)
		}

		return out
	} else {
		return nil
	}
}

// Returns the first item that is not the zero value for that value's type.
func Or(in ...interface{}) interface{} {
	for _, v := range Compact(in) {
		// if the current value equals the zero value of its type,
		// then skip it, otherwise return it
		if v != reflect.Zero(reflect.TypeOf(v)).Interface() {
			return v
		}
	}

	return nil
}

// Returns the first item that is not a zero-length string.
func OrString(in ...string) string {
	if v := CompactString(in); len(v) > 0 {
		return v[0]
	} else {
		return ``
	}
}

// Returns the length of the given slice, array, or string.
func Len(in interface{}) int {
	in = typeutil.ResolveValue(in)

	if typeutil.IsKind(in, reflect.Array, reflect.Slice, reflect.String) {
		inV := reflect.ValueOf(in)
		return inV.Len()
	}

	return 0
}

// Returns the element in the given indexable value at the given index.  If the
// index is present, the second return value will be true.  If the index is not
// present, or the given input is not indexable, the second return value will be
// false.
func At(in interface{}, index int) (interface{}, bool) {
	in = typeutil.ResolveValue(in)

	if typeutil.IsKind(in, reflect.Array, reflect.Slice, reflect.String) {
		inV := reflect.ValueOf(in)

		if index < inV.Len() {
			return inV.Index(index).Interface(), true
		}
	}

	return nil, false
}

// Returns the nth element from the given slice, array or string; or nil.
func Get(in interface{}, index int) interface{} {
	if v, ok := At(in, index); ok {
		return v
	} else {
		return nil
	}
}

// Returns the first element from the given slice, array or string; or nil.
func First(in interface{}) interface{} {
	return Get(in, 0)
}

// Returns the first element in the given inputs that is not that type's zero value.  All input values
// are flattened into a single array, so variadic elements can contain scalar or array values.
func FirstNonZero(inputs ...interface{}) interface{} {
	for _, v := range Flatten(inputs) {
		if !typeutil.IsZero(v) {
			return v
		}
	}

	return nil
}

// Returns the all but the first element from the given slice, array or string; or nil.
func Rest(in interface{}) []interface{} {
	if typeutil.IsKind(in, reflect.Array, reflect.Slice, reflect.String) {
		inV := reflect.ValueOf(in)
		l := inV.Len()

		switch l {
		case 0, 1:
			return nil
		default:
			out := make([]interface{}, l-1)

			for i := 1; i < l; i++ {
				elemV := inV.Index(i)

				if elemV.CanInterface() {
					out[i-1] = elemV.Interface()
				}
			}

			return out
		}
	}

	return nil
}

// Returns the last element from the given slice, array or string; or nil.
func Last(in interface{}) interface{} {
	if Len(in) == 0 {
		return nil
	}

	return Get(in, Len(in)-1)
}

// Iterate through each element of the given array or slice, calling
// iterFn exactly once for each element.  Otherwise, call iterFn one time
// with the given input as the argument.
func Each(slice interface{}, iterFn IterationFunc) error {
	return utils.SliceEach(slice, iterFn)
}

// Takes some input value and returns it as a slice.
func Sliceify(in interface{}) []interface{} {
	return utils.Sliceify(in)
}

// Returns a new slice with only the specified subset of items included.  In addition to the
// normal slice index rules in Golang, negative indices are also supported.  If a negative index is
// given for the from and/or to values, the index will be treated as being relative to the end of the
// given slice. For example:
//
//	Slice([]interface{}{1,2,3,4,5}, -5, -1)  // returns []interface{}{1, 2, 3, 4, 5}
//	Slice([]interface{}{1,2,3,4,5}, -2, -1)  // returns []interface{}{4, 5}
//	Slice([]interface{}{1,2,3,4,5}, -1, -1)  // returns []interface{}{5}
//	Slice([]interface{}{1,2,3,4,5}, -4, -2)  // returns []interface{}{2, 3, 4}
func Slice(slice interface{}, from int, to int) []interface{} {
	sliceS := Sliceify(slice)

	if from < 0 {
		from = len(sliceS) + from
	}

	if from > len(sliceS) {
		return make([]interface{}, 0)
	} else if from < 0 {
		from = 0
	}

	if to > len(sliceS) {
		to = len(sliceS)
	} else if to < 0 {
		to = len(sliceS) + to + 1
	}

	if (from >= 0 && from < len(sliceS)) && (to >= from && to <= len(sliceS)) {
		return sliceS[from:to]
	} else {
		return make([]interface{}, 0)
	}
}

// Same as slice, but returns strings.
func StringSlice(slice interface{}, from int, to int) []string {
	return Stringify(Slice(slice, from, to))
}

// Returns a new slice with only unique elements from the given interface included.
func Unique(in interface{}) []interface{} {
	return unique(in, StrictEqualityCompare)
}

// Returns a new slice with only unique string elements from the given interface included.
func UniqueStrings(in interface{}) []string {
	return Stringify(Unique(in))
}

func unique(in interface{}, comparator CompareFunc) []interface{} {
	inV := reflect.ValueOf(in)
	values := make([]interface{}, 0)

	if inV.IsValid() {
	InputLoop:
		for i := 0; i < inV.Len(); i++ {
			if current := inV.Index(i); current.IsValid() {
				for _, existing := range values {
					if comparator(i, existing, current.Interface()) {
						continue InputLoop
					}
				}

				values = append(values, current.Interface())
			}
		}
	} else {
		return nil
	}

	return values
}

// Returns a copy of the given slice with each element modified by the a given function.
func Map(in interface{}, fn MapFunc) []interface{} {
	var out []interface{}

	Each(in, func(i int, v interface{}) error {
		out = append(out, fn(i, v))
		return nil
	})

	return out
}

// Returns a copy of the given slice with each element modified by the a given function, then
// converted to a string.
func MapString(in interface{}, fn MapStringFunc) []string {
	out := Stringify(in)

	for i, el := range out {
		out[i] = fn(i, el)
	}

	return out
}

// Trims the whitespace from each element in the given array.
func TrimSpace(in interface{}) []string {
	return MapString(in, func(_ int, el string) string {
		return strings.TrimSpace(el)
	})
}

// Divide the given slice into chunks of (at most) a given length
func Chunks(in interface{}, size int) [][]interface{} {
	if !typeutil.IsArray(in) {
		return nil
	}

	output := make([][]interface{}, 0)
	current := make([]interface{}, 0)

	Each(in, func(i int, v interface{}) error {
		if i > 0 && i%size == 0 {
			output = append(output, current)
			current = nil
		}

		current = append(current, v)
		return nil
	})

	if len(current) > 0 {
		output = append(output, current)
	}

	return output
}

// Returns a copy of the given slice with each element's value passed to stringutil.Autotype
func Autotype(in interface{}) []interface{} {
	var out []interface{}

	Each(in, func(i int, v interface{}) error {
		out = append(out, stringutil.Autotype(v))
		return nil
	})

	return out
}

// Returns a copy of the given slicified value with the given additional values appended.
func Append(in interface{}, values ...interface{}) []interface{} {
	return append(Sliceify(in), values...)
}
