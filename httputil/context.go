package httputil

import (
	"context"
	"net/http"

	"github.com/ghetzel/go-stockutil/typeutil"
)

// attach an arbitrary value to the context of a given request.
func RequestSetValue(req *http.Request, key string, value any) {
	var parent = req.Context()
	var withValue = context.WithValue(parent, key, value)
	*req = *req.WithContext(withValue)
}

// Retrieve an arbitrary value from the context of a given request.
func RequestGetValue(req *http.Request, key string) typeutil.Variant {
	if value := req.Context().Value(key); value != nil {
		return typeutil.V(value)
	} else {
		return typeutil.V(nil)
	}
}
