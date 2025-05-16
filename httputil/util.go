package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

// Encode the username and password into a value than can be used in the Authorization HTTP header.
func EncodeBasicAuth(username string, password string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", username, password)),
	))
}

// Configures the given http.Client to accept TLS certificates validated by the given PEM-encoded CA bundle file
func SetRootCABundle(client *http.Client, caBundle string) error {
	return updateRootCABundle(false, client, caBundle)
}

// Loads certificates from the given file and returns a usable x509.CertPool
func LoadCertPool(filename string) (*x509.CertPool, error) {
	if data, err := fileutil.ReadAll(filename); err == nil {
		var pool = x509.NewCertPool()

		if pool.AppendCertsFromPEM(data) {
			return pool, nil
		} else {
			return nil, fmt.Errorf("An error occurred adding the provided certificate(s)")
		}
	} else {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}
}

// Returns the media type from the request's Content-Type.
func MediaType(req *http.Request) string {
	var contentType = req.Header.Get(`Content-Type`)

	if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
		return mediaType
	} else {
		return contentType
	}
}

// Returns whether the Content-Type of the given request matches any of the supplied options.
// The mediaTypes arguments may be either complete media types (e.g.: "text/html", "application/javascript")
// or major type classes (e.g.: "text/", "video/").  The trailing slash (/) indicates that any media type
// that begins with that text will match.
func IsMediaType(req *http.Request, mediaTypes ...string) bool {
	var mediaType = MediaType(req)

	for _, mt := range mediaTypes {
		if strings.HasSuffix(mt, `/`) {
			if strings.HasPrefix(mediaType, mt) {
				return true
			}
		} else if mt == mediaType {
			return true
		}
	}

	return false
}

// UrlPathJoin takes a string or *url.URL and joins the existing URL path component with the given path.
// The new path may also contain query string values, which will be added to the base URL.  Existing keys will
// be replaced with new ones, except for repeated keys (e.g.: ?x=1&x=2&x=3).  In this case, the new values will
// be added to the existing ones.  The *url.URL returned from this function is a copy, and the original URL (if
// one is provided) will not be modified in any way.
func UrlPathJoin(baseurl any, path string) (*url.URL, error) {
	var in *url.URL
	var out *url.URL

	if u, ok := baseurl.(*url.URL); ok {
		in = u
	} else if u, err := url.Parse(typeutil.String(baseurl)); err == nil {
		in = u
	} else {
		return nil, err
	}

	newpath, qs := stringutil.SplitPair(path, `?`)
	var trail string

	if strings.HasSuffix(newpath, `/`) {
		newpath = strings.TrimSuffix(newpath, `/`)
		in.Path = strings.TrimSuffix(in.Path, `/`)
		trail = `/`
	}

	out = new(url.URL)
	out.Scheme = in.Scheme
	out.Opaque = in.Opaque
	out.User = in.User
	out.Host = in.Host
	out.Path = filepath.Join(in.Path, newpath) + trail
	out.RawPath = in.RawPath
	out.ForceQuery = in.ForceQuery
	out.RawQuery = in.RawQuery
	out.Fragment = in.Fragment

	if qs != `` {
		if qsv, err := url.ParseQuery(qs); err == nil {
			for k, vs := range qsv {
				if len(vs) == 1 {
					SetQ(out, k, vs[0])
				} else {
					AddQ(out, k, sliceutil.Sliceify(vs)...)
				}
			}
		} else {
			return nil, err
		}
	}

	return out, nil
}

func updateRootCABundle(appendPem bool, client *http.Client, caBundle string) error {
	if caBundle == `` {
		return fmt.Errorf("Must specify a file to read certificates from")
	} else if client == nil {
		return fmt.Errorf("Must provide an *http.Client to modify")
	} else {
		if pool, err := LoadCertPool(caBundle); err == nil {
			if client.Transport == nil {
				client.Transport = &http.Transport{}
			}

			if htt, ok := client.Transport.(*http.Transport); ok {
				if htt.TLSClientConfig == nil {
					htt.TLSClientConfig = &tls.Config{}
				}

				htt.TLSClientConfig.RootCAs = pool
				return nil
			} else {
				return fmt.Errorf("Cannot configure TLS on HTTP transport %T", client.Transport)
			}
		} else {
			return err
		}
	}
}
