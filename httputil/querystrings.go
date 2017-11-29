package httputil

import (
	"net/http"
	"net/url"
	"time"

	"github.com/ghetzel/go-stockutil/stringutil"
)

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

func QTime(req *http.Request, key string) time.Time {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToTime(v); err == nil {
			return i
		}
	}

	return time.Time{}
}

func QBool(req *http.Request, key string) bool {
	if v, err := stringutil.ConvertToBool(Q(req, key)); err == nil {
		return v
	}

	return false
}

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