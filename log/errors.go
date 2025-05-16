package log

import (
	"strings"

	"github.com/ghetzel/go-stockutil/typeutil"
)

func precheck(err error, message any) (string, string, bool) {
	if err == nil || message == nil {
		return ``, ``, false
	}

	if me, ok := message.(error); ok {
		return err.Error(), me.Error(), true
	} else {
		return err.Error(), typeutil.String(message), true
	}
}

// Return whether the given error is prefixed with the given message.  Message can
// be a string or another error.  If either is nil, this function returns false.
func ErrHasPrefix(err error, message any) bool {
	if emsg, msg, ok := precheck(err, message); ok {
		return strings.HasPrefix(emsg, msg)
	} else {
		return false
	}
}

// Return whether the given error contains with the given message.  Message can
// be a string or another error.  If either is nil, this function returns false.
func ErrContains(err error, message any) bool {
	if emsg, msg, ok := precheck(err, message); ok {
		return strings.Contains(emsg, msg)
	} else {
		return false
	}
}

// Return whether the given error is suffixed with the given message.  Message can
// be a string or another error.  If either is nil, this function returns false.
func ErrHasSuffix(err error, message any) bool {
	if emsg, msg, ok := precheck(err, message); ok {
		return strings.HasSuffix(emsg, msg)
	} else {
		return false
	}
}
