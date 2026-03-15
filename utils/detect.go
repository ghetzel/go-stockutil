package utils

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func IsInteger(in any) bool {
	var inV = reflect.ValueOf(in)

	switch inV.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true

	default:
		if asStr, err := ToString(in); err == nil {
			if _, err := strconv.Atoi(asStr); err == nil {
				return true
			}
		}
	}

	return false
}

func IsFloat(in any) bool {
	var inV = reflect.ValueOf(in)

	switch inV.Kind() {
	case reflect.Float32, reflect.Float64:
		return true

	default:
		if asStr, err := ToString(in); err == nil {
			if _, err := strconv.ParseFloat(asStr, 64); err == nil {
				return true
			}
		}
	}

	return false
}

func IsHexadecimal(in any) bool {
	if inS, err := ToString(in); err == nil {
		inS = strings.ToLower(inS)

		if after, ok := strings.CutPrefix(inS, `0x`); ok {
			inS = after
		} else {
			return false
		}

		for _, r := range inS {
			if r >= '0' && r <= '9' || r >= 'a' && r <= 'f' {
				continue
			} else {
				return false
			}
		}
	} else {
		return false
	}

	return true
}

func IsNumeric(in any) bool {
	return IsFloat(in)
}

func IsBoolean(inI any) bool {
	if in, err := ToString(inI); err == nil {
		in = strings.ToLower(in)

		return (IsBooleanTrue(in) || IsBooleanFalse(in))
	}

	return false
}

func IsBooleanTrue(inI any) bool {
	if in, err := ToString(inI); err == nil {
		in = strings.ToLower(in)

		if slices.Contains(BooleanTrueValues, in) {
			return true
		}
	}

	return false
}

func IsBooleanFalse(inI any) bool {
	if in, err := ToString(inI); err == nil {
		in = strings.ToLower(in)

		if slices.Contains(BooleanFalseValues, in) {
			return true
		}
	}

	return false
}

func IsTime(inI any) bool {
	if in, err := ToString(inI); err == nil {
		if f := DetectTimeFormat(in); f != `` && f != `epoch` {
			return true
		}
	}

	return false
}
