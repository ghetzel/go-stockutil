package httputil

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

// Deprecated: this type will go away in 1.9.x
type RequestLogger struct {
}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

func (self *RequestLogger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var start = time.Now()

	next(rw, req)

	var response = rw.(negroni.ResponseWriter)
	var status = response.Status()
	var duration = time.Since(start)

	if status < 400 {
		Logger.Debugf("[HTTP %d] %s to %v took %v", status, req.Method, req.URL, duration)
	} else {
		Logger.Warningf("[HTTP %d] %s to %v took %v", status, req.Method, req.URL, duration)
	}
}
