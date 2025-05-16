package httputil

import (
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dsnet/compress/brotli"
)

// Takes an http.Response and returns an io.Reader that will return the contents
// of the Response Body decoded according to the values (if any) of the Content-Encoding
// response header.
func DecodeResponse(response *http.Response) (io.Reader, error) {
	var output io.Reader = response.Body

	var contentEncodings = strings.Split(response.Header.Get(`Content-Encoding`), `,`)

	for _, enc := range contentEncodings {
		enc = strings.TrimSpace(enc)
		enc = strings.TrimPrefix(enc, `x-`)

		if c, err := decode(output, enc); err == nil {
			output = c
		} else {
			return nil, err
		}
	}

	return output, nil
}

func decode(input io.Reader, encoding string) (io.Reader, error) {
	switch encoding {
	case `identity`, ``:
		return input, nil

	case `gzip`:
		return gzip.NewReader(input)

	case `deflate`:
		return flate.NewReader(input), nil

	case `bzip2`:
		return bzip2.NewReader(input), nil

	case `br`:
		return brotli.NewReader(input, nil)

	default:
		return nil, fmt.Errorf("Unsupported encoding %q", encoding)
	}
}
