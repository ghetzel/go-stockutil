package utils

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ghetzel/go-stockutil/rxutil"
)

var rxExtendedDurations = regexp.MustCompile(`(?i)((?P<number>[\d\.]+)(?P<suffix>[^\d]+))`)

const durationMaxMatches = 16

type ConvertType int

const (
	Invalid ConvertType = iota
	Bytes
	String
	Float
	Integer
	Time
	Boolean
	Nil
	UserDefined
)

func (self ConvertType) String() string {
	switch self {
	case String:
		return `str`
	case Boolean:
		return `bool`
	case Float:
		return `float`
	case Integer:
		return `int`
	case Time:
		return `time`
	case Bytes:
		return `bytes`
	case UserDefined:
		return `user`
	default:
		return ``
	}
}

func (self ConvertType) IsSupersetOf(other ConvertType) bool {
	return self < other
}

var PassthroughType = errors.New(`passthrough`)

type TypeConvertFunc func(in any) (any, error)

var typeHandlers = make(map[string]TypeConvertFunc)

// Register's a handler used for converting one type to another. Type are checked in the following
// manner:  The input value's reflect.Type String() value is matched, falling back to its
// reflect.Kind String() value, finally checking for a special "*" value that matches any type.
// If the handler function returns nil, its value replaces the input value.  If the special error
// type PassthroughType is returned, the original value is returned unmodified.
func RegisterTypeHandler(handler TypeConvertFunc, types ...string) {
	for _, t := range types {
		if t != `` {
			typeHandlers[t] = handler
		}
	}
}

var rxLeadingZeroes = regexp.MustCompile(`^0+\d+$`)
var NilStrings = []string{`null`, `NULL`, `<nil>`, `nil`, `Nil`, `None`, `undefined`, ``}
var BooleanTrueValues = []string{`true`, `yes`, `on`}
var BooleanFalseValues = []string{`false`, `no`, `off`}

var TimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC850,
	time.RFC822,
	time.RFC822Z,
	time.RFC1123,
	time.RFC1123Z,
	time.Kitchen,
	`2006-01-02 15:04:05.000000000`,
	`2006-01-02 15:04:05.000000`,
	`2006-01-02 15:04:05.000`,
	`2006-01-02 15:04:05 -0700 MST`,
	`2006-01-02 15:04:05Z07:00`,
	`2006-01-02 15:04:05`,
	`2006-01-02 15:04`,
	`2006-01-02`,
	`2006-01-02T15:04:05.000000000`,
	`2006-01-02T15:04:05.000000`,
	`2006-01-02T15:04:05.000`,
	`2006-01-02T15:04:05 -0700 MST`,
	`2006-01-02T15:04:05Z07:00`,
	`2006-01-02T15:04:05`,
	`2006-01-02T15:04`,
}

func ToString(in any) (string, error) {
	if in == nil {
		return ``, nil
	} else if err, ok := in.(error); ok {
		return err.Error(), nil
	} else if s, ok := in.(fmt.Stringer); ok {
		return s.String(), nil
	}

	var asBytes []byte

	if u8, ok := in.([]uint8); ok {
		asBytes = []byte(u8)
	} else if b, ok := in.([]byte); ok {
		asBytes = b
	} else if r, ok := in.([]rune); ok {
		return string(r), nil
	}

	if len(asBytes) > 0 {
		if out := string(asBytes); utf8.ValidString(out) {
			return out, nil
		} else {
			return ``, fmt.Errorf("Given %T is not a valid UTF-8 string", in)
		}
	}

	if inT := reflect.TypeOf(in); inT != nil {
		switch inT.Kind() {
		case reflect.Float32:
			return strconv.FormatFloat(reflect.ValueOf(in).Float(), 'f', -1, 32), nil
		case reflect.Float64:
			return strconv.FormatFloat(reflect.ValueOf(in).Float(), 'f', -1, 64), nil
		case reflect.Bool:
			return strconv.FormatBool(in.(bool)), nil
		case reflect.String:
			if inStr, ok := in.(string); ok {
				return inStr, nil
			}
		}

		if !IsKind(in, CompoundTypes...) {
			return fmt.Sprintf("%v", in), nil
		}
	}

	return ``, fmt.Errorf("Unable to convert type '%T' to string", in)
}

func ConvertTo(toType ConvertType, inI any) (any, error) {
	var inS string
	var inSerr error

	if inV, ok := inI.(reflect.Value); ok {
		if inV.CanInterface() {
			inI = inV.Interface()
		} else {
			return nil, fmt.Errorf("reflect.Value given, but cannot retrieve interface value")
		}
	}

	inS, inSerr = ToString(inI)

	switch toType {
	case Float:
		if inS == `` {
			return float64(0), nil
		} else if v, ok := inI.(float32); ok {
			return float64(v), nil
		} else if v, ok := inI.(float64); ok {
			return float64(v), nil
		} else if v, ok := inI.(int); ok {
			return float64(v), nil
		} else if v, ok := inI.(int8); ok {
			return float64(v), nil
		} else if v, ok := inI.(int16); ok {
			return float64(v), nil
		} else if v, ok := inI.(int32); ok {
			return float64(v), nil
		} else if v, ok := inI.(int64); ok {
			return float64(v), nil
		} else if v, ok := inI.(uint); ok {
			return float64(v), nil
		} else if v, ok := inI.(uint8); ok {
			return float64(v), nil
		} else if v, ok := inI.(uint16); ok {
			return float64(v), nil
		} else if v, ok := inI.(uint32); ok {
			return float64(v), nil
		} else if v, ok := inI.(uint64); ok {
			return float64(v), nil
		} else if IsHexadecimal(inI) {
			v, err := ConvertHexToInteger(inI)
			return float64(v), err
		}

		return strconv.ParseFloat(inS, 64)
	case Integer:
		if inTime, ok := inI.(time.Time); ok {
			return int64(inTime.UnixNano()), nil
		} else if inS == `` {
			return int64(0), nil
		} else if layout := DetectTimeFormat(inS); layout != `` && layout != `epoch` {
			if tm, err := time.Parse(layout, inS); err == nil {
				return tm.UnixNano(), nil
			} else {
				return nil, err
			}
		} else if IsHexadecimal(inI) {
			return ConvertHexToInteger(inI)
		}

		if v, ok := inI.(int); ok {
			return int64(v), nil
		} else if v, ok := inI.(int8); ok {
			return int64(v), nil
		} else if v, ok := inI.(int16); ok {
			return int64(v), nil
		} else if v, ok := inI.(int32); ok {
			return int64(v), nil
		} else if v, ok := inI.(int64); ok {
			return int64(v), nil
		} else if v, ok := inI.(uint); ok {
			return int64(v), nil
		} else if v, ok := inI.(uint8); ok {
			return int64(v), nil
		} else if v, ok := inI.(uint16); ok {
			return int64(v), nil
		} else if v, ok := inI.(uint32); ok {
			return int64(v), nil
		} else if v, ok := inI.(uint64); ok {
			return int64(v), nil
		}

		return strconv.ParseInt(inS, 10, 64)
	case Boolean:
		if inI == nil {
			return false, nil
		}

		if IsBooleanTrue(inS) {
			return true, nil
		} else if IsBooleanFalse(inS) {
			return false, nil
		} else {
			return nil, fmt.Errorf("Cannot convert '%s' into a boolean value", inS)
		}
	case Time:
		if inTime, ok := inI.(time.Time); ok {
			return inTime, nil
		}

		inS = strings.Trim(strings.TrimSpace(inS), `"'`)

		if DetectTimeFormat(inS) == `epoch` {
			if v, err := strconv.ParseInt(inS, 10, 64); err == nil {
				return time.Unix(v, 0), nil
			}
		}

		for _, format := range TimeFormats {
			if tm, err := time.Parse(format, strings.TrimSpace(inS)); err == nil {
				return tm, nil
			}
		}

		switch inS {
		case `now`:
			return time.Now(), nil
		default:
			// handle time zero values
			var tmS = strings.Map(func(r rune) rune {
				switch r {
				case '-', ':', ' ', 'T', 'Z':
					return '0'
				}

				return r
			}, inS)

			if v, err := strconv.ParseInt(tmS, 10, 64); err == nil && v == 0 {
				return time.Time{}, nil
			}

			return nil, fmt.Errorf("Cannot convert '%s' into a date/time value", inS)
		}

	case Bytes:
		if inI == nil {
			return []byte{}, nil
		} else if inR, ok := inI.(io.Reader); ok {
			// special case: read all io.Reader, convert resulting bytes to string
			if data, err := io.ReadAll(inR); err == nil {
				return data, nil
			} else {
				return ``, fmt.Errorf("Cannot convert io.Reader to []byte: %v", err)
			}
		} else if inB, ok := inI.([]byte); ok {
			return inB, nil
		} else if inB, ok := inI.([]uint8); ok {
			var outB = make([]byte, len(inB))

			for i, v := range inB {
				outB[i] = byte(v)
			}

			return outB, nil
		} else if IsKind(inI, reflect.Slice, reflect.Array) {
			var outB = make([]byte, 0)

			if err := SliceEach(inI, func(i int, value any) error {
				if bb, err := ConvertToInteger(value); err == nil {
					outB = append(outB, byte(bb))
					return nil
				} else {
					return err
				}
			}); err == nil {
				return outB, nil
			} else {
				return nil, err
			}
		} else if inSerr == nil {
			return []byte(inS), inSerr
		} else {
			return nil, inSerr
		}

	case String:
		if inI == nil {
			return ``, nil
		} else if inR, ok := inI.(io.Reader); ok {
			// special case: read all io.Reader, convert resulting bytes to string
			if data, err := io.ReadAll(inR); err == nil {
				return string(data), nil
			} else {
				return ``, fmt.Errorf("Cannot convert io.Reader to string: %v", err)
			}
		} else if inB, ok := inI.([]byte); ok {
			// special case: assume incoming byte slices are actually strings
			// convert byte slices to strings directly
			return string(inB), nil
		} else if inSr, ok := inI.(fmt.Stringer); ok {
			// convert fmt.Stringer to string
			return inSr.String(), nil

		} else if inV := reflect.ValueOf(inI); inV.Kind() == reflect.Ptr {
			// dereference pointers to strings and stringify the result
			inS, inSerr = ToString(inV.Elem())
		}

		return inS, inSerr

	default:
		return inI, nil
	}
}

func ConvertHexToInteger(in any) (int64, error) {
	if IsHexadecimal(in) {
		if inS, err := ToString(in); err == nil {
			inS = strings.ToLower(inS)
			inS = strings.TrimPrefix(inS, `0x`)

			return strconv.ParseInt(inS, 16, 64)
		} else {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("invalid hexadecimal value '%v'", in)
	}
}

func ConvertToInteger(in any) (int64, error) {
	if v, err := ConvertTo(Integer, in); err == nil {
		return v.(int64), nil
	} else {
		return int64(0), err
	}
}

func ConvertToFloat(in any) (float64, error) {
	if v, err := ConvertTo(Float, in); err == nil {
		return v.(float64), nil
	} else {
		return float64(0.0), err
	}
}

func ConvertToString(in any) (string, error) {
	if v, err := ConvertTo(String, in); err == nil {
		return v.(string), nil
	} else {
		return ``, err
	}
}

func ConvertToBool(in any) (bool, error) {
	if v, err := ConvertTo(Boolean, in); err == nil {
		return v.(bool), nil
	} else {
		return false, err
	}
}

func ConvertToTime(in any) (time.Time, error) {
	switch in.(type) {
	case time.Time:
		return in.(time.Time), nil
	default:
		if v, err := ConvertTo(Time, in); err == nil {
			return v.(time.Time), nil
		} else {
			return time.Time{}, err
		}
	}
}

func ConvertToBytes(in any) ([]byte, error) {
	if v, err := ConvertTo(Bytes, in); err == nil {
		return v.([]byte), nil
	} else {
		return nil, err
	}
}

func DetectTimeFormat(in string) string {
	if IsInteger(in) {
		return `epoch`
	}

	for _, layout := range TimeFormats {
		if _, err := time.Parse(layout, in); err == nil {
			return layout
		}
	}

	return ``
}

// Returns the given value, converted according to any handlers set via RegisterTypeHandler.
func ConvertCustomType(in any) (any, error) {
	var convert TypeConvertFunc
	var inV reflect.Value

	// if we were given a reflect.Value target, then we shouldn't take the reflect.ValueOf that
	if tV, ok := in.(reflect.Value); ok {
		inV = tV
	} else {
		inV = reflect.ValueOf(in)
	}

	if inV.IsValid() {
		// give type handlers a chance to do their work
		if handler, ok := typeHandlers[inV.Type().String()]; ok {
			convert = handler
		} else if handler, ok := typeHandlers[inV.Kind().String()]; ok {
			convert = handler
		} else if handler, ok := typeHandlers[`*`]; ok {
			convert = handler
		}

		if convert != nil {
			if out, err := convert(in); err == nil {
				return out, nil
			} else {
				return nil, err
			}
		}
	}

	return in, PassthroughType
}

func Detect(in any) (ConvertType, any) {
	// perform custom type conversions (if any)
	if v, err := ConvertCustomType(in); err == nil {
		return UserDefined, v
	} else if err != PassthroughType {
		panic(err.Error())
	}

	if in == nil {
		return Nil, nil
	}

	if IsTime(in) {
		if v, err := ConvertTo(Time, in); err == nil {
			return Time, v
		}
	}

	// effectively, this detects strings that are numeric, but have leading zeroes.
	// we should treat those as meaningful and return that string outright
	//
	// (e.g.: handle the "US Zip Code" problem; i.e.: 07753 _can_ be interpreted as int(7753))
	if vStr, ok := in.(string); ok {
		if rxLeadingZeroes.MatchString(vStr) {
			return String, vStr
		}

		// certain known string values should convert to nil directly
		for _, nilStr := range NilStrings {
			if vStr == nilStr {
				return Nil, nil
			}
		}
	}

	for _, ctype := range []ConvertType{
		Boolean,
		Integer,
		Float,
		String,
	} {
		if value, err := ConvertTo(ctype, in); err == nil {
			return ctype, value
		}
	}

	return Invalid, in
}

func Autotype(in any) any {
	_, value := Detect(in)
	return value
}

func DetectConvertType(in any) ConvertType {
	ctype, _ := Detect(in)
	return ctype
}

func ParseDuration(in string) (time.Duration, error) {
	if in == `` {
		return 0, nil
	}

	var i int
	var totalHours int

	in = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, in)

	for {
		i++

		if match := rxutil.Match(rxExtendedDurations, in); match != nil {
			if num, err := ConvertToInteger(match.Group(`number`)); err == nil {
				var hours int

				switch strings.ToLower(match.Group(`suffix`)) {
				case `year`, `years`, `y`:
					hours = 8760
				case `week`, `weeks`, `wk`, `wks`, `w`:
					hours = 168
				case `day`, `days`, `d`:
					hours = 24
				case `hour`, `hours`, `hr`, `hrs`, `h`:
					hours = 1
				case `minute`, `minutes`, `min`:
					in = match.ReplaceGroup(`suffix`, `m`)
				default:
					break
				}

				if hours > 0 {
					in = match.ReplaceGroup(1, ``)
				}

				totalHours += int(num) * hours
			} else {
				return 0, fmt.Errorf("Invalid number: %v", err)
			}
		} else {
			break
		}

		if i >= durationMaxMatches {
			break
		}
	}

	if totalHours > 0 {
		in = fmt.Sprintf("%dh%s", totalHours, in)
	}

	return time.ParseDuration(in)
}
