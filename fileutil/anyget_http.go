package fileutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var HttpDefaultTimeout = 10 * time.Second

// Retrieve a file via HTTP or HTTPS.
//
// Supported Context Values:
//
//	insecure:
//	  (bool) specify that strict TLS validation should be optional.
//
//	method:
//	  (string) the HTTP method to use, defaults to GET.
//
//	metadata:
//	  (map[string]interface{}) a key-value set of HTTP request headers to include.
//
//	safeResponseCodes:
//	  ([]int) a list of one or more HTTP status codes that are considered successful for this request.
func RetrieveViaHTTP(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	ctx, timeout := ctxToTimeout(ctx, HttpDefaultTimeout)

	var client = &http.Client{
		Timeout: timeout,
	}

	if typeutil.Bool(ctx.Value(`insecure`)) {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if req, err := http.NewRequest(
		typeutil.OrString(ctx.Value(`method`), http.MethodGet),
		u.String(),
		nil,
	); err == nil {
		for kv := range maputil.M(ctx.Value(`metadata`)).Iter() {
			if k := kv.K; k != `` {
				if v := kv.V.String(); v != `` {
					req.Header.Set(k, v)
				}
			}
		}

		var statusOk = func(res *http.Response) bool {
			var codes = sliceutil.Stringify(ctx.Value(`safeResponseCodes`))

			if len(codes) > 0 {
				return sliceutil.ContainsString(codes, typeutil.String(res.StatusCode))
			} else {
				return res.StatusCode < 400
			}
		}

		if res, err := client.Do(req); err == nil {
			if statusOk(res) {
				if res.Body != nil {
					return res.Body, nil
				} else {
					return NullReadCloser, nil
				}
			} else {
				return nil, fmt.Errorf("responded HTTP: %v", res.Status)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
