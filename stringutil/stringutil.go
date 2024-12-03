// Utilities for converting and manipulating data to and from strings
package stringutil

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/go-stockutil/utils"
	"github.com/jdkato/prose/tokenize"
)

var rxHexadecimal = regexp.MustCompile(`^[0-9a-fA-F]+$`)
var rxEnvVarExpr = regexp.MustCompile(`(\$\{(?P<env>\w+)(?:\|(?P<fallback>.*?))?(?::(?P<fmt>%[^\}]+))?\})`) // ${ENV}, ${ENV:%04s}, ${ENV|fallback}, ${ENV|fallback:%04s}
var DefaultThousandsSeparator = `,`
var DefaultDecimalSeparator = `.`

// If set to true, ExpandEnv() will preserve ${...} sequences whose resulting value would yield an empty string.
// Fallback and formatting values are still supported, as the determination of value emptiness is made after parsing
// the sequence.  This is to avoid unset environment variables resulting in the escape sequences being erased, which is often
// the case when running ExpandEnv() against various shell languages (Bash, et. al)
var ExpandEnvPreserveIfEmpty = false
var ExpandEnvTempDelimiterOpen = "\u3018"  // LEFT WHITE TORTOISE SHELL BRACKET (U+3018, Ps): 〘
var ExpandEnvTempDelimiterClose = "\u3019" // RIGHT WHITE TORTOISE SHELL BRACKET (U+3019, Pe): 〙

var NilStrings = utils.NilStrings
var BooleanTrueValues = utils.BooleanTrueValues
var BooleanFalseValues = utils.BooleanFalseValues
var TimeFormats = utils.TimeFormats

type SiPrefix int

const (
	None  SiPrefix = 0
	Kilo           = 1
	Mega           = 2
	Giga           = 3
	Tera           = 4
	Peta           = 5
	Exa            = 6
	Zetta          = 7
	Yotta          = 8
)

func (self SiPrefix) String() string {
	switch self {
	case Kilo:
		return `K`
	case Mega:
		return `M`
	case Giga:
		return `G`
	case Tera:
		return `T`
	case Peta:
		return `P`
	case Exa:
		return `E`
	case Zetta:
		return `Z`
	case Yotta:
		return `Y`
	default:
		return ``
	}
}

type ConvertType = utils.ConvertType

const (
	Invalid ConvertType = utils.Invalid
	Nil                 = utils.Nil
	String              = utils.String
	Boolean             = utils.Boolean
	Float               = utils.Float
	Integer             = utils.Integer
	Time                = utils.Time
	Bytes               = utils.Bytes
)

func ParseType(name string) ConvertType {
	switch strings.ToLower(name) {
	case `str`:
		return utils.String
	case `bool`:
		return utils.Boolean
	case `float`:
		return utils.Float
	case `int`:
		return utils.Integer
	case `time`:
		return utils.Time
	case `bytes`:
		return utils.Bytes
	default:
		return utils.Invalid
	}
}

func IsInteger(in interface{}) bool {
	return utils.IsInteger(in)
}

func IsFloat(in interface{}) bool {
	return utils.IsFloat(in)
}

func IsNumeric(in interface{}) bool {
	return utils.IsNumeric(in)
}

func IsBoolean(in interface{}) bool {
	return utils.IsBoolean(in)
}

func IsBooleanTrue(in interface{}) bool {
	return utils.IsBooleanTrue(in)
}

func IsBooleanFalse(in interface{}) bool {
	return utils.IsBooleanFalse(in)
}

func IsTime(in interface{}) bool {
	return utils.IsTime(in)
}

func IsSurroundedBy(inI interface{}, prefix string, suffix string) bool {
	if in, err := ToString(inI); err == nil {
		if strings.HasPrefix(in, prefix) && strings.HasSuffix(in, suffix) {
			return true
		}
	}

	return false
}

func DetectTimeFormat(in string) string {
	return utils.DetectTimeFormat(in)
}

func ToString(in interface{}) (string, error) {
	return utils.ToString(in)
}

func MustString(in interface{}, fallbackOpt ...string) string {
	if v, err := ToString(in); err == nil {
		return v
	} else if len(fallbackOpt) > 0 {
		return fallbackOpt[0]
	} else {
		panic(err.Error())
	}
}

func MustInteger(in interface{}, fallbackOpt ...int64) int64 {
	if v, err := ConvertToInteger(in); err == nil {
		return v
	} else if len(fallbackOpt) > 0 {
		return fallbackOpt[0]
	} else {
		panic(err.Error())
	}
}

func MustFloat(in interface{}, fallbackOpt ...float64) float64 {
	if v, err := ConvertToFloat(in); err == nil {
		return v
	} else if len(fallbackOpt) > 0 {
		return fallbackOpt[0]
	} else {
		panic(err.Error())
	}
}

func MustBool(in interface{}, fallbackOpt ...bool) bool {
	if v, err := ConvertToBool(in); err == nil {
		return v
	} else if len(fallbackOpt) > 0 {
		return fallbackOpt[0]
	} else {
		panic(err.Error())
	}
}

func MustTime(in interface{}, fallbackOpt ...time.Time) time.Time {
	if v, err := ConvertToTime(in); err == nil {
		return v
	} else if len(fallbackOpt) > 0 {
		return fallbackOpt[0]
	} else {
		panic(err.Error())
	}
}

func ToStringSlice(in interface{}) ([]string, error) {
	values := make([]string, 0)

	if in != nil {
		if v, ok := in.([]string); ok {
			return v, nil
		}

		inV := reflect.ValueOf(in)

		if inV.IsValid() {
			if inV.Kind() == reflect.Ptr {
				inV = inV.Elem()
			}

			if inV.IsValid() {
				switch inV.Kind() {
				case reflect.Array, reflect.Slice:
					for i := 0; i < inV.Len(); i++ {
						if indexV := inV.Index(i); indexV.IsValid() {
							if v, err := ToString(indexV.Interface()); err == nil {
								values = append(values, v)
							} else {
								return nil, err
							}
						} else {
							return nil, fmt.Errorf("Element %d in slice is invalid", i)
						}
					}

				default:
					if v, err := ToString(in); err == nil {
						values = append(values, v)
					} else {
						return nil, err
					}
				}
			} else {
				return nil, fmt.Errorf("Cannot parse value pointed to by given value.")
			}

		} else {
			return nil, fmt.Errorf("Cannot parse given value.")
		}
	}

	return values, nil
}

func ToByteString(in interface{}, formatString ...string) (string, error) {
	if asBytes, err := ConvertToInteger(in); err == nil {
		for i := 0; i < 9; i++ {
			if converted := (float64(asBytes) / math.Pow(1024, float64(i))); converted < 1024 {
				prefix := SiPrefix(i)
				f := `%g`

				if len(formatString) > 0 {
					f = formatString[0]
				}

				return fmt.Sprintf(f+"%sB", converted, prefix.String()), nil
			}
		}

		return fmt.Sprintf("%dB", asBytes), nil
	} else {
		return ``, err
	}
}

func GetSiPrefix(input string) (SiPrefix, error) {
	switch input {
	case "", "b", "B":
		return None, nil
	case "k", "K":
		return Kilo, nil
	case "m", "M":
		return Mega, nil
	case "g", "G":
		return Giga, nil
	case "t", "T":
		return Tera, nil
	case "p", "P":
		return Peta, nil
	case "e", "E":
		return Exa, nil
	case "z", "Z":
		return Zetta, nil
	case "y", "Y":
		return Yotta, nil
	default:
		return None, fmt.Errorf("Unrecognized SI unit '%s'", input)
	}
}

func ToBytes(input string) (float64, error) {
	//  handle -ibibyte suffixes like KiB, GiB
	if strings.HasSuffix(input, "ib") || strings.HasSuffix(input, "iB") {
		input = input[0 : len(input)-2]

		//  handle input that puts the 'B' in the suffix; e.g.: Kb, GB
	} else if len(input) > 2 && IsInteger(string(input[len(input)-3])) && (input[len(input)-1] == 'b' || input[len(input)-1] == 'B') {
		input = input[0 : len(input)-1]
	}

	if prefix, err := GetSiPrefix(string(input[len(input)-1])); err == nil {
		if v, err := strconv.ParseFloat(input[0:len(input)-1], 64); err == nil {
			return v * math.Pow(1024, float64(prefix)), nil
		} else {
			return 0, err
		}
	} else {
		if v, err := strconv.ParseFloat(input, 64); err == nil {
			return v, nil
		} else {
			return 0, fmt.Errorf("Unrecognized input string '%s'", input)
		}
	}
}

func ConvertTo(toType ConvertType, inI interface{}) (interface{}, error) {
	return utils.ConvertTo(toType, inI)
}

func ConvertToInteger(in interface{}) (int64, error) {
	return utils.ConvertToInteger(in)
}

func ConvertToFloat(in interface{}) (float64, error) {
	return utils.ConvertToFloat(in)
}

func ConvertToString(in interface{}) (string, error) {
	return utils.ConvertToString(in)
}

func ConvertToBool(in interface{}) (bool, error) {
	return utils.ConvertToBool(in)
}

func ConvertToTime(in interface{}) (time.Time, error) {
	return utils.ConvertToTime(in)
}

func ConvertToBytes(in interface{}) ([]byte, error) {
	return utils.ConvertToBytes(in)
}

func Autotype(in interface{}) interface{} {
	return utils.Autotype(in)
}

func IsSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		}

		return true
	}

	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}

	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

func TokenizeFunc(in string, tokenizer func(rune) bool, partfn func(part string) []string) []string {
	// split on word-separating characters (and discard them), or on capital
	// letters (preserving them)
	parts := strings.FieldsFunc(in, tokenizer)
	out := make([]string, 0)

	for _, part := range parts {
		partOut := partfn(part)

		if partOut != nil {
			for _, v := range partOut {
				if v != `` {
					out = append(out, v)
				}
			}
		}
	}

	return out
}

func Camelize(in interface{}) string {
	return strings.Join(TokenizeFunc(MustString(in), IsSeparator, func(part string) []string {
		part = strings.TrimSpace(part)
		part = strings.Title(part)
		return []string{part}
	}), ``)
}

func Underscore(in interface{}) string {
	return Snakeify(in, '_')
}

func Hyphenate(in interface{}) string {
	return Snakeify(in, '-')
}

func Snakeify(in interface{}, separator rune) string {
	inS := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return separator
		} else if unicode.IsPunct(r) && separator != r {
			return separator
		} else {
			return r
		}
	}, MustString(in))

	out := make([]rune, 0)
	runes := []rune(inS)

	sepfn := func(i int) bool {
		return i >= 0 && i < len(runes) && unicode.IsLower(runes[i])
	}

	for i, r := range runes {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)

			if i > 0 && runes[i-1] != separator && (sepfn(i-1) || sepfn(i+1)) {
				out = append(out, separator)
			}
		}

		out = append(out, r)
	}

	return string(out)
}

// Returns whether the letters (Unicode Catgeory 'L') in a given string are
// homogenous in case (all upper-case or all lower-case).
func IsMixedCase(in string) bool {
	var hasLower bool
	var hasUpper bool

	for _, c := range in {
		if unicode.IsLetter(c) {
			if unicode.IsLower(c) {
				hasLower = true

				if hasUpper {
					return true
				}
			} else if unicode.IsUpper(c) {
				hasUpper = true

				if hasLower {
					return true
				}
			}
		}
	}

	return false
}

// Returns whether the given string is a hexadecimal number. If the string is
// prefixed with "0x", the prefix is removed first. If length is greater than 0,
// the length of the input (excluding prefix) is checked as well.
func IsHexadecimal(in string, length int) bool {
	in = strings.TrimPrefix(in, `0x`)

	if IsMixedCase(in) {
		return false
	}

	if rxHexadecimal.MatchString(in) {
		if length <= 0 {
			return true
		} else if len(in) == length {
			return true
		}
	}

	return false
}

func Thousandify(in interface{}, separator string, decimal string) string {
	if separator == `` {
		separator = DefaultThousandsSeparator
	}

	if decimal == `` {
		decimal = DefaultDecimalSeparator
	}

	if inStr, err := ToString(in); err == nil {
		if IsNumeric(in) {
			var buffer []rune

			lastIndexBeforeDecimal := strings.Index(inStr, decimal) - 1
			decimalAndAfter := strings.Index(inStr, decimal)

			if lastIndexBeforeDecimal < 0 {
				lastIndexBeforeDecimal = len(inStr) - 1
			}

			j := 0

			for i := lastIndexBeforeDecimal; i >= 0; i-- {
				j++
				buffer = append([]rune{rune(inStr[i])}, buffer...)

				if j == 3 && i > 0 && !(i == 1 && inStr[0] == '-') {
					buffer = append([]rune(separator), buffer...)
					j = 0
				}
			}

			if decimalAndAfter >= 0 {
				for _, r := range inStr[decimalAndAfter:] {
					buffer = append(buffer, rune(r))
				}
			}

			return string(buffer[:])
		} else {
			return inStr
		}
	} else {
		return ``
	}
}

func LongestCommonPrefix(inputs []string) string {
	output := ``
	shortestInputLen := 0

	for _, in := range inputs {
		if shortestInputLen == 0 || len(in) < shortestInputLen {
			shortestInputLen = len(in)
		}
	}

LCPLoop:
	for i := 0; i < shortestInputLen; i++ {
		var current byte

		for _, input := range inputs {
			if i < len(input) {
				if current == 0 {
					current = input[i]
				} else if current != input[i] {
					break LCPLoop
				}
			}
		}

		if current == 0 {
			break
		}

		output = output + string(current)
	}

	return output
}

func RelaxedEqual(first interface{}, second interface{}) (bool, error) {
	if reflect.DeepEqual(first, second) {
		return true, nil
	} else if IsNumeric(first) && IsNumeric(second) {
		if fV, err := ConvertToFloat(first); err == nil {
			if sV, err := ConvertToFloat(second); err == nil {
				return (fV == sV), nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	} else if IsBooleanTrue(first) && IsBooleanTrue(second) {
		return true, nil
	} else if IsBooleanFalse(first) && IsBooleanFalse(second) {
		return true, nil
	} else if IsTime(first) && IsTime(second) {
		if fV, err := ConvertToTime(first); err == nil {
			if sV, err := ConvertToTime(second); err == nil {
				return fV.Equal(sV), nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	} else {
		if fV, err := ToString(first); err == nil {
			if sV, err := ToString(second); err == nil {
				return (fV == sV), nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}
}

// A case insensitive, single return version of RelaxedEqual that trims leading and trailing whitespace from strings before comparison.
func SoftEqual(first interface{}, second interface{}) bool {
	if typeutil.IsKindOfString(first) {
		first = strings.TrimSpace(strings.ToLower(typeutil.String(first)))
	}

	if typeutil.IsKindOfString(second) {
		second = strings.TrimSpace(strings.ToLower(typeutil.String(second)))
	}

	if eq, err := RelaxedEqual(first, second); err == nil && eq {
		return true
	} else {
		return false
	}
}

// Split the given string into two parts.  If there is only one resulting part,
// that part will be the first return value and the second return value will be empty.
func SplitPair(in string, delimiter string) (string, string) {
	parts := strings.Split(in, delimiter)

	switch len(parts) {
	case 1:
		return parts[0], ``
	default:
		return parts[0], strings.Join(parts[1:], delimiter)
	}
}

// Identical to SplitPair, but trims any leading and trailing whitespace from each output value.
func SplitPairTrimSpace(in string, delimiter string) (string, string) {
	var a, b = SplitPair(in, delimiter)

	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	return a, b
}

// Identical to SplitPair, but returns the result of passing the second value in the pair to Autotype()
func SplitPairAuto(in string, delimiter string) (string, interface{}) {
	var a, b = SplitPair(in, delimiter)

	return a, Autotype(b)
}

// Identical to SplitPairTrimSpace, but returns the result of passing the second value in the pair to Autotype()
func SplitPairTrimSpaceAuto(in string, delimiter string) (string, interface{}) {
	var a, b = SplitPairTrimSpace(in, delimiter)

	return a, Autotype(b)
}

// Split the given string into two parts.  If there is only one resulting part,
// that part will be the second return value and the first return value will be empty.
func SplitPairTrailing(in string, delimiter string) (string, string) {
	first, rest := SplitPair(in, delimiter)

	if rest == `` {
		return rest, first
	} else {
		return first, rest
	}
}

// Split the given string into two parts from the right. If there is only one resulting part,
// that part will be the first return value and the second return value will be empty.
func SplitPairRight(in string, delimiter string) (string, string) {
	parts := strings.Split(in, delimiter)

	switch len(parts) {
	case 1:
		return parts[0], ``
	default:
		return strings.Join(parts[0:len(parts)-1], delimiter), parts[len(parts)-1]
	}
}

// Split the given string into two parts.  If there is only one resulting part,
// that part will be the second return value and the first return value will be empty.
func SplitPairRightTrailing(in string, delimiter string) (string, string) {
	first, rest := SplitPairRight(in, delimiter)

	if rest == `` {
		return rest, first
	} else {
		return first, rest
	}
}

func SplitTriple(in string, delimiter string) (string, string, string) {
	parts := strings.Split(in, delimiter)

	switch len(parts) {
	case 1:
		return parts[0], ``, ``
	case 2:
		return parts[0], parts[1], ``
	default:
		return parts[0], parts[1], strings.Join(parts[2:], delimiter)
	}
}

// Identical to SplitPair, but trims any leading and trailing whitespace from each output value.
func SplitTripleTrimSpace(in string, delimiter string) (string, string, string) {
	var a, b, c = SplitTriple(in, delimiter)

	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	c = strings.TrimSpace(c)

	return a, b, c
}

// Identical to strings.Split, but trims any leading and trailing whitespace from each element.
func SplitTrimSpace(s string, sep string) []string {
	var parts = strings.Split(s, sep)

	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts
}

// Identical to strings.Split, but trims any leading and trailing whitespace from each element using
// strings.Trim.
func SplitTrim(s string, sep string, cutset string) []string {
	var parts = strings.Split(s, sep)

	for i, part := range parts {
		parts[i] = strings.Trim(part, cutset)
	}

	return parts
}

// Prefix the given string if it is non-empty
func PrefixIf(in string, prefix string) string {
	if !typeutil.IsZero(in) {
		in = prefix + in
	}

	return in
}

// Suffix the given string if it is non-empty
func SuffixIf(in string, suffix string) string {
	if !typeutil.IsZero(in) {
		in = in + suffix
	}

	return in
}

// Return the given string with prefixed and suffixed with other strings.
func Wrap(in string, prefix string, suffix string) string {
	in = strings.TrimPrefix(in, prefix)
	in = strings.TrimSuffix(in, suffix)

	return prefix + in + suffix
}

// Return the given string with the given prefix and suffix removed.
func Unwrap(in string, prefix string, suffix string) string {
	in = strings.TrimPrefix(in, prefix)
	in = strings.TrimSuffix(in, suffix)
	return in
}

// Wrap the given string if it is non-empty
func WrapIf(in string, prefix string, suffix string) string {
	if !typeutil.IsZero(in) {
		in = Wrap(in, prefix, suffix)
	}

	return in
}

// Wrap each element in the given string slice with prefix and suffix.
func WrapEach(in []string, prefix string, suffix string) []string {
	out := make([]string, len(in))

	for i, v := range in {
		out[i] = Wrap(v, prefix, suffix)
	}

	return out
}

// Prefix each element in the given string slice with prefix.
func PrefixEach(in []string, prefix string) []string {
	return WrapEach(in, prefix, ``)
}

// Suffix each element in the given string slice with suffix.
func SuffixEach(in []string, prefix string, suffix string) []string {
	return WrapEach(in, ``, suffix)
}

// Split the given input into lines.
func SplitLines(in interface{}, delimiter string) []string {
	var blob string

	if b, ok := in.([]byte); ok {
		blob = string(b)
	} else {
		blob = typeutil.String(in)
	}

	return strings.Split(blob, delimiter)
}

// Split the given string into words.
func SplitWords(in string) []string {
	tokenizer := tokenize.NewTreebankWordTokenizer()
	out := make([]string, 0)

	for _, word := range tokenizer.Tokenize(in) {
		out = append(out, word)
	}

	return out
}

// Truncate the given string to a certain number of characters.
func Elide(in string, charcount int, trailer ...string) string {
	if len(in) < charcount {
		return in
	}

	return in[0:charcount] + strings.Join(trailer, ``)
}

// Truncate the given string to a certain number of characters from the end.
func ElideRight(in string, charcount int, leader ...string) string {
	if len(in) < charcount {
		return in
	}

	return strings.Join(leader, ``) + in[len(in)-charcount:]
}

// Truncate the given string to a certain number of words.
func ElideWords(in string, wordcount int) string {
	words := SplitWords(in)
	wc := len(words)

	if wc == 0 {
		return ``
	} else if wc <= wordcount {
		return strings.Join(words, ` `)
	} else {
		words = words[0:wordcount]

		return strings.TrimRightFunc(strings.Join(words, ` `), func(r rune) bool {
			return unicode.IsPunct(r) || unicode.IsSpace(r)
		})
	}
}

type OnlySqueezeFunc func(r rune) bool

// Return the given string with sequences of characters matching the given function
// replaced with a single instance of that character.
func SqueezeFunc(in string, fn OnlySqueezeFunc) string {
	out := make([]rune, 0)
	var previous rune

	for _, char := range in {
		if fn == nil || fn(char) {
			if char != previous {
				out = append(out, char)
			}
		} else {
			out = append(out, char)
		}

		previous = char
	}

	return string(out)
}

// Return the given string with sequences repeating character replaced with a
// single instance of that character.
func Squeeze(in string) string {
	return SqueezeFunc(in, func(r rune) bool {
		return true
	})
}

// Return the given string with sequences of whitespace characters replaced with
// a single instance of that character.
func SqueezeSpace(in string) string {
	return SqueezeFunc(in, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

// Return the given string with environment variable substitution sequences
// expanded and (optionally) formatted.  This function operates similarly to
// os.ExpandEnv, but accepts custom fmt.Printf formatting directives.
func ExpandEnv(in string) string {
	for {
		if match := rxutil.Match(rxEnvVarExpr, in); match != nil {
			var format = match.Group(`fmt`)
			var varname = match.Group(`env`)
			var fallback = match.Group(`fallback`)

			if varname != `` {
				if format == `` {
					format = `%v`
				}

				var typed interface{}

				if val := os.Getenv(varname); val != `` {
					typed = Autotype(val)
				} else {
					typed = Autotype(fallback)
				}

				if typed != nil {
					in = match.ReplaceGroup(1, fmt.Sprintf(format, typed))
				} else if ExpandEnvPreserveIfEmpty {
					// this One Weird Trick™ always weirds me out, but here goes
					// we're in a loop that recursively expands all ${...} sequences and eliminates those sequences
					// from the output *UNLESS* ExpandEnvPreserveIfEmpty==true.  In that case, we're saying
					// "if this sequence would expand to an empty string, leave it alone in the output."
					//
					// Problem is, that creates an infinite loop of replacing the same sequence with itself forever.
					// So what we do instead is we swap out the delimiters denoting the start and end of the expansion
					// sequence with a Unicode character that is _very_ unlikely to appear in the types of strings this
					// function is expected to process.  We then go back, outside of the loop, and swap those characters back
					// with the correct delimiters.
					//
					// It's a bit janky, for sure, but it works. If you somehow end up with a conflict with these specific
					// temporary delimiters, you can change the package-level variables to some less-problematic characters
					// to your use case.

					in = match.ReplaceGroup(1, fmt.Sprintf(format, ExpandEnvTempDelimiterOpen+varname+ExpandEnvTempDelimiterClose))
				} else {
					in = match.ReplaceGroup(1, fmt.Sprintf(format, ``))
				}
			}
		} else {
			break
		}
	}

	in = strings.ReplaceAll(in, ExpandEnvTempDelimiterOpen, `${`)
	in = strings.ReplaceAll(in, ExpandEnvTempDelimiterClose, `}`)

	return in
}

// Takes the given string, splits it into lines, and prefixes each line with the given prefix string.
func PrefixLines(in interface{}, prefix string) string {
	lines := SplitLines(in, "\n")

	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}

// Takes the given string, splits it into lines, and suffixes each line with the given suffix string.
func SuffixLines(in interface{}, suffix string) string {
	lines := SplitLines(in, "\n")

	for i, line := range lines {
		lines[i] = line + suffix
	}

	return strings.Join(lines, "\n")
}

// HasPrefix tests whether the string s begins with any prefix in prefixes.
func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}

	return false
}

// HasSuffix tests whether the string s begins with any suffix in suffixes.
func HasAnySuffix(s string, suffixes ...string) bool {
	for _, p := range suffixes {
		if strings.HasSuffix(s, p) {
			return true
		}
	}

	return false
}
