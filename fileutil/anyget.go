package fileutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
)

type RetrieveHandlerFunc = func(context.Context, *url.URL) (io.ReadCloser, error)

var GenericDefaultTimeout = 5 * time.Second
var NullReadCloser = io.NopCloser(bytes.NewBuffer(nil))
var retrieveSchemeHandlers = make(map[string]RetrieveHandlerFunc)

// Register a new Retrieve() function for the given URL scheme.
func RegisterRetrieveScheme(scheme string, fn RetrieveHandlerFunc) {
	retrieveSchemeHandlers[scheme] = fn
}

// Unregister the given URL scheme from being handled by the Retrieve() function.
func UnregisterRetrieveScheme(scheme string) {
	delete(retrieveSchemeHandlers, scheme)
}

// Perform a generic retrieval of data located at a specified resource given as a URL.
// This function supports file://, http://, https://, ssh://, and sftp:// schemes, and can
// be extended to support additional schemes using the RegisterRetrieveScheme package function.
//
// If resourceUri is given as a *url.URL, the value of that URL will be copied.  Any other type
// will be converted to a string (honoring types that implement fmt.Stringer), and the resulting
// URL will be used.
func Retrieve(ctx context.Context, resourceUri any) (io.ReadCloser, error) {
	ctx, _ = ctxToTimeout(ctx, 0)

	var uri url.URL

	if u, ok := resourceUri.(*url.URL); ok {
		uri = *u
	} else {
		var r = typeutil.String(resourceUri)

		if !strings.Contains(r, `://`) {
			r = `file:///` + r
		}

		if u, err := url.Parse(r); err == nil {
			uri = *u
		} else {
			return nil, fmt.Errorf("bad url: %v", err)
		}
	}

	uri.Scheme = strings.ToLower(uri.Scheme)

	if handler, ok := retrieveSchemeHandlers[uri.Scheme]; ok && handler != nil {
		return handler(ctx, &uri)
	} else {
		return nil, fmt.Errorf("unsupported scheme %q", uri.Scheme)
	}
}

func ctxToTimeout(ctx context.Context, fallback time.Duration) (context.Context, time.Duration) {
	var timeout time.Duration

	if ctx == nil {
		ctx = context.Background()
	}

	if dl, ok := ctx.Deadline(); ok {
		timeout = time.Until(dl)
	} else if fallback > 0 {
		timeout = fallback
	} else {
		timeout = GenericDefaultTimeout
	}

	return ctx, timeout
}

func init() {
	RegisterRetrieveScheme(`file`, RetrieveViaFilesystem)
	RegisterRetrieveScheme(`http`, RetrieveViaHTTP)
	RegisterRetrieveScheme(`https`, RetrieveViaHTTP)
	RegisterRetrieveScheme(`ssh`, RetrieveViaSSH)
	RegisterRetrieveScheme(`sftp`, RetrieveViaSSH)
}
