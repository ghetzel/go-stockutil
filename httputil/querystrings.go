package httputil

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/timeutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

// Parses the named query string from a request as an integer.
func QInt(req *http.Request, key string, fallbacks ...int64) int64 {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToInteger(v); err == nil {
			return i
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return 0
	}
}

// Parses the named query string from a request as a float.
func QFloat(req *http.Request, key string, fallbacks ...float64) float64 {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToFloat(v); err == nil {
			return i
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return 0
	}
}

// Parses the named query string from a request as a date/time value.
func QTime(req *http.Request, key string) time.Time {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToTime(v); err == nil {
			return i
		}
	}

	return time.Time{}
}

// Parses the named query string from a request as a duration string.
func QDuration(req *http.Request, key string) time.Duration {
	if v := Q(req, key); v != `` {
		if d, err := timeutil.ParseDuration(v); err == nil {
			return d
		}
	}

	return 0
}

// Parses the named query string from a request as a boolean value.
func QBool(req *http.Request, key string, fallbacks ...bool) bool {
	if v := Q(req, key); v == `` && len(fallbacks) > 0 {
		return fallbacks[0]
	} else if typeutil.Bool(v) {
		return true
	}

	return false
}

// Parses the named query string from a request as a string.
func Q(req *http.Request, key string, fallbacks ...string) string {
	if v := req.URL.Query().Get(key); v != `` {
		if vS, err := url.QueryUnescape(v); err == nil {
			return vS
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return ``
	}
}

// Parses the named query string from a request as a delimiter-separated string slice.
func QStrings(req *http.Request, key string, delimiter string, fallbacks ...string) []string {
	if strs := sliceutil.CompactString(strings.Split(Q(req, key), delimiter)); len(strs) > 0 {
		return strs
	} else if len(fallbacks) > 0 {
		return sliceutil.Stringify(sliceutil.Flatten(fallbacks))
	} else {
		return make([]string, 0)
	}
}

// Sets a query string to the given value in the given url.URL
func SetQ(u *url.URL, key string, value any) {
	var qs = u.Query()
	qs.Set(key, stringutil.MustString(value))
	u.RawQuery = qs.Encode()
}

// Appends a query string from then given url.URL
func AddQ(u *url.URL, key string, values ...any) {
	var qs = u.Query()

	for _, value := range values {
		qs.Add(key, stringutil.MustString(value))
	}

	u.RawQuery = qs.Encode()
}

// Deletes a query string from then given url.URL
func DelQ(u *url.URL, key string) {
	var qs = u.Query()
	qs.Del(key)
	u.RawQuery = qs.Encode()
}

// A version of SetQ that accepts a URL string and makes a best-effort to modify it.
// Will return the modified URL or the original URL if an error occurred.
func SetQString(u string, key string, value any) string {
	if ur, err := url.Parse(u); err == nil {
		SetQ(ur, key, value)

		return ur.String()
	}

	return u
}

// A version of AddQ that accepts a URL string and makes a best-effort to modify it.
// Will return the modified URL or the original URL if an error occurred.
func AddQString(u string, key string, value any) string {
	if ur, err := url.Parse(u); err == nil {
		AddQ(ur, key, value)

		return ur.String()
	}

	return u
}

// A version of DelQ that accepts a URL string and makes a best-effort to modify it.
// Will return the modified URL or the original URL if an error occurred.
func DelQString(u string, key string) string {
	if ur, err := url.Parse(u); err == nil {
		DelQ(ur, key)

		return ur.String()
	}

	return u
}
