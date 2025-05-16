package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
	multierror "github.com/hashicorp/go-multierror"
	"k8s.io/client-go/util/jsonpath"
)

var ReferenceTime time.Time = time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.FixedZone("MST", -7*60*60))

func GenericMarshalJSON(in any, extraData ...map[string]any) ([]byte, error) {
	var sval = structs.New(in)
	var output = make(map[string]any)

FieldLoop:
	for _, field := range sval.Fields() {
		if field.IsExported() {
			var key = field.Name()
			var value = field.Value()

			// perform type conversions because encoding/json won't Do The Right Thing in all cases.
			if vTime, ok := value.(time.Time); ok {
				if vTime.IsZero() {
					value = nil
				} else {
					value = vTime.Format(time.RFC3339Nano)
				}
			} else if _, ok := value.(fmt.Stringer); ok {
				value = fmt.Sprintf("%v", value)
			}

			// parse "json" tag
			if parts := strings.Split(field.Tag(`json`), `,`); len(parts) > 0 {
				if parts[0] != `` {
					key = parts[0]
				}

				if len(parts) > 1 {
					for _, flag := range parts[1:] {
						switch flag {
						case `omitempty`:
							if field.IsZero() {
								continue FieldLoop
							}
						}
					}
				}
			}

			output[key] = value
		}
	}

	if len(extraData) > 0 && len(extraData[0]) > 0 {
		for k, v := range extraData[0] {
			output[k] = v
		}
	}

	return json.Marshal(output)
}

// Appends on error to another, allowing for operations that return multiple errors
// to remain compatible within a single-valued context.
func AppendError(base error, err error) error {
	if err == nil {
		return base
	} else {
		return multierror.Append(base, err)
	}
}

// Performs a JSONPath query against the given object and returns the results.
// JSONPath description, syntax, and examples are available at http://goessner.net/articles/JsonPath/.
func JSONPath(data any, query string, autowrap bool) (any, error) {
	if reflect.TypeOf(data).Kind() == reflect.Map && query != `` {
		for _, line := range strings.Split(query, "\n") {
			line = strings.TrimSpace(line)

			if line == `` {
				continue
			} else if autowrap {
				line = strings.TrimPrefix(line, `{`)
				line = strings.TrimSuffix(line, `}`)
				line = `{` + line + `}`
			}

			var jp = jsonpath.New(``).AllowMissingKeys(true)

			if err := jp.Parse(line); err == nil {
				var values []any

				if results, err := jp.FindResults(data); err == nil {
					for _, pair := range results {
						for _, p := range pair {
							if p.IsValid() && p.CanInterface() {
								values = append(values, p.Interface())
							}
						}
					}
				} else {
					return nil, err
				}

				switch len(values) {
				case 0:
					data = nil
				case 1:
					data = values[0]
				default:
					data = values
				}
			} else {
				return nil, err
			}
		}
	}

	return data, nil
}
