package httputil

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	gohttputil "net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var DebugOutputBoxWidth = 60
var WaitForPollInterval = time.Second

var DefaultClient = MustClient(``)

// A simplified GET function using the package-level default client.  Will return the
// response body as bytes.
func GetBody(url string) ([]byte, error) {
	if res, err := DefaultClient.Get(url, nil, nil); err == nil {
		var data, err = io.ReadAll(res.Body)

		res.Body.Close()

		return data, err
	} else {
		return nil, err
	}
}

type Literal []byte

// Periodically performs a GET request against the given URL, waiting up to timeout
// for a 200-series HTTP response code.
func WaitForHTTP(url string, timeout time.Duration, c ...*http.Client) error {
	var client *http.Client

	if len(c) > 0 && c[0] != nil {
		client = c[0]
	} else {
		client = http.DefaultClient
	}

	var start = time.Now()

	for time.Since(start) < timeout {
		if res, err := client.Get(url); err == nil {
			if res.StatusCode < 400 {
				return nil
			}
		}

		time.Sleep(WaitForPollInterval)
	}

	return fmt.Errorf("Request to %s did not succeed in %v", url, timeout)
}

type Client struct {
	encoder           EncoderFunc
	decoder           DecoderFunc
	errorDecoder      ErrorDecoderFunc
	firstRequestHook  InitRequestFunc
	preRequestHook    InterceptRequestFunc
	postRequestHook   InterceptResponseFunc
	uri               *url.URL
	headers           map[string]any
	params            map[string]any
	httpClient        *http.Client
	rootCaPool        *x509.CertPool
	insecure          bool
	firstRequestSent  bool
	basicAuthViaNetrc bool
}

func MustClient(baseURI string) *Client {
	if c, err := NewClient(baseURI); err == nil {
		return c
	} else {
		panic(err.Error())
	}
}

func NewClient(baseURI string) (*Client, error) {
	var client = &Client{
		encoder:    JSONEncoder,
		decoder:    JSONDecoder,
		headers:    make(map[string]any),
		params:     make(map[string]any),
		httpClient: http.DefaultClient,
	}

	if baseURI != `` {
		if uri, err := url.Parse(baseURI); err == nil {
			client.uri = uri
		} else {
			return nil, err
		}
	} else {
		client.uri = new(url.URL)
	}

	if log.VeryDebugging(`github.com/ghetzel/go-stockutil/httputil`) {
		client.SetPreRequestHook(func(req *http.Request) (any, error) {
			if data, err := gohttputil.DumpRequest(req, true); err == nil {
				log.Debugf("httputil ${blue}\u256d\u2500[ HTTP Request ]%s\u2504\u2504\u2504${reset}", strings.Repeat("\u2500", DebugOutputBoxWidth-17))
				log.Debugf("httputil ${blue}\u2502${reset} \u21c9 %v", req.URL)
				log.Debugf("httputil ${blue}\u2502${reset}")

				for _, line := range strings.Split(string(data), "\n") {
					line = strings.TrimSpace(line)
					log.Debugf("httputil ${blue}\u2502${reset} %v", line)
				}

				log.Debugf("httputil ${blue}\u2570%s\u2504\u2504\u2504${reset}", strings.Repeat("\u2500", DebugOutputBoxWidth))
			}

			return nil, nil
		})

		client.SetPostRequestHook(func(res *http.Response, _ any) error {
			if data, err := gohttputil.DumpResponse(res, true); err == nil {
				log.Debugf("httputil \u256d\u2500[ HTTP Response ]%s\u2504\u2504\u2504", strings.Repeat("\u2500", DebugOutputBoxWidth-18))

				if res.Request != nil {
					log.Debugf("httputil ${red}\u2502${reset} \u21c7 %v", res.Request.URL)
					log.Debugf("httputil ${red}\u2502${reset}")
				}

				for _, line := range strings.Split(string(data), "\n") {
					line = strings.TrimSpace(line)
					log.Debugf("httputil \u2502 %v", line)
				}

				log.Debugf("httputil \u2570%s\u2504\u2504\u2504", strings.Repeat("\u2500", DebugOutputBoxWidth))
			}

			return nil
		})
	}

	return client, nil
}

// Return a copy of the current client that uses a different encoder.
func (self *Client) WithEncoder(fn EncoderFunc) *Client {
	var clientWith = new(Client)
	*clientWith = *self
	clientWith.encoder = fn
	return clientWith
}

// Return a copy of the current client that uses a different decoder.
func (self *Client) WithDecoder(fn DecoderFunc) *Client {
	var clientWith = new(Client)
	*clientWith = *self
	clientWith.decoder = fn
	return clientWith
}

// Return the base URI for this client.
func (self *Client) URI() *url.URL {
	return self.uri
}

// Specify that the NetrcFile (default: ~/.netrc) should be consulted before each
// request to supply basic authentication.  If a non-empty username or password is
// found for the
func (self *Client) SetAutomaticLogin(on bool) {
	self.basicAuthViaNetrc = on
}

// Specify an encoder that will be used to serialize data in the request body.
func (self *Client) SetEncoder(fn EncoderFunc) {
	self.encoder = fn
}

// Specify a decoder that will be used to deserialize data in the response body.
func (self *Client) SetDecoder(fn DecoderFunc) {
	self.decoder = fn
}

// Specify a different decoder used to deserialize non 2xx/3xx HTTP responses.
func (self *Client) SetErrorDecoder(fn ErrorDecoderFunc) {
	self.errorDecoder = fn
}

// Specify a function that will be called immediately before the first request is sent.
// This function has an opportunity to read and modify the outgoing request, and
// if it returns a non-nil error, the request will not be sent.
func (self *Client) SetInitHook(fn InitRequestFunc) {
	self.firstRequestHook = fn
}

// Specify a function that will be called immediately before a request is sent.
// This function has an opportunity to read and modify the outgoing request, and
// if it returns a non-nil error, the request will not be sent.
func (self *Client) SetPreRequestHook(fn InterceptRequestFunc) {
	self.preRequestHook = fn
}

// Specify a function tht will be called immediately after a response is received.
// This function is given the first opportunity to inspect the response, and if it
// returns a non-nil error, no additional processing (including the Error Decoder function)
// will be performed.
func (self *Client) SetPostRequestHook(fn InterceptResponseFunc) {
	self.postRequestHook = fn
}

// Return the headers set on this client.
func (self *Client) Headers() http.Header {
	var hdr = make(http.Header)

	for k, v := range self.headers {
		hdr.Add(k, typeutil.String(v))
	}

	return hdr
}

// Remove all implicit HTTP request headers.
func (self *Client) ClearHeaders() {
	self.headers = nil
}

// Add an HTTP request header by name that will be included in every request. If
// value is nil, the named header will be removed instead.
func (self *Client) SetHeader(name string, value any) {
	if value != nil {
		self.headers[name] = value
	} else {
		delete(self.headers, name)
	}
}

// Remove all implicit querystring parameters.
func (self *Client) ClearParams() {
	self.params = nil
}

// Add a querystring parameter by name that will be included in every request. If
// value is nil, the parameter will be removed instead.
func (self *Client) SetParam(name string, value any) {
	if value != nil {
		self.params[name] = value
	} else {
		delete(self.params, name)
	}
}

// Return the params set on this client.
func (self *Client) Params() map[string]any {
	return self.params
}

// Returns the HTTP client used to perform requests
func (self *Client) Client() *http.Client {
	return self.httpClient
}

// Replace the default HTTP client with a user-provided one
func (self *Client) SetClient(client *http.Client) {
	self.httpClient = client
}

// Set or unset insecure TLS requests that will proceed even if the peer certificate cannot be verified.
func (self *Client) SetInsecureTLS(insecure bool) {
	defer self.syncHttpClient()
	self.insecure = insecure
}

// Set the username and password to be included in the Authorization header.
func (self *Client) SetBasicAuth(username string, password string) {
	self.SetHeader(`Authorization`, EncodeBasicAuth(username, password))
}

// Append one or more trusted certificates to the RootCA bundle that is consulted when performing HTTPS requests.
func (self *Client) AppendTrustedRootCA(pemFilenamesOrData ...any) error {
	return self.updateRootCA(false, pemFilenamesOrData...)
}

// Replace the existing RootCA bundle with an explicit set of trusted certificates.
func (self *Client) SetRootCA(pemFilenamesOrData ...any) error {
	return self.updateRootCA(true, pemFilenamesOrData...)
}

func (self *Client) updateRootCA(replace bool, pemFilenamesOrData ...any) error {
	var pems = make([][]byte, 0)

	// when we're all done, make sure the http.Client we'll be using knows about our certs
	defer self.syncHttpClient()

	for _, pem := range pemFilenamesOrData {
		if b, ok := pem.([]byte); ok {
			pems = append(pems, b)
		} else if r, ok := pem.(io.Reader); ok {
			if b, err := io.ReadAll(r); err == nil {
				pems = append(pems, b)
			} else {
				return err
			}
		} else if s, ok := pem.(string); ok {
			if b, err := fileutil.ReadAll(s); err == nil {
				pems = append(pems, b)
			} else {
				return err
			}
		} else {
			return fmt.Errorf("type error: expected []byte, io.Reader, or string; got %T", pem)
		}
	}

	// clear out or setup the initial RootCA pool
	if replace {
		self.rootCaPool = x509.NewCertPool()
	} else if self.rootCaPool == nil {
		// try to inherit the RootCA from any existing client
		if transport, ok := self.httpClient.Transport.(*http.Transport); ok {
			if tcc := transport.TLSClientConfig; tcc != nil {
				self.rootCaPool = tcc.RootCAs
			}
		}

		// still nil? get the existing system CA bundle
		if self.rootCaPool == nil {
			if syspool, err := x509.SystemCertPool(); err == nil {
				self.rootCaPool = syspool
			} else {
				return fmt.Errorf("Failed to retrieve system CA pool: %v", err)
			}
		}
	}

	// append each cert to the root pool
	for _, pem := range pems {
		if !self.rootCaPool.AppendCertsFromPEM(pem) {
			return fmt.Errorf("Failed to append certificate")
		}
	}

	return nil
}

func (self *Client) syncHttpClient() {
	var newTCC = &tls.Config{
		RootCAs:            self.rootCaPool,
		InsecureSkipVerify: self.insecure,
	}

	newTCC.BuildNameToCertificate()

	// handle use cases where an existing client is present
	if transport, ok := self.httpClient.Transport.(*http.Transport); ok {
		if tcc := transport.TLSClientConfig; tcc != nil {
			tcc.RootCAs = newTCC.RootCAs
			tcc.InsecureSkipVerify = self.insecure
		} else {
			transport.TLSClientConfig = newTCC
		}
	} else {
		self.httpClient.Transport = &http.Transport{
			TLSClientConfig: newTCC,
		}
	}
}

// Perform an HTTP request
func (self *Client) Request(
	method Method,
	path string,
	body any,
	params map[string]any,
	headers map[string]any,
) (*http.Response, error) {
	return self.RequestWithContext(context.Background(), method, path, body, params, headers)
}

// Perform an HTTP request using the given context
func (self *Client) RequestWithContext(
	ctx context.Context,
	method Method,
	path string,
	body any,
	params map[string]any,
	headers map[string]any,
) (*http.Response, error) {
	if !self.firstRequestSent {
		if self.firstRequestHook != nil {
			if err := self.firstRequestHook(); err != nil {
				return nil, fmt.Errorf("init hook: %v", err)
			}
		}

		self.firstRequestSent = true
	}

	var finalParams = make(map[string]any)
	var finalHeaders = make(map[string]any)

	for k, v := range self.params {
		finalParams[k] = v
	}

	for k, v := range params {
		finalParams[k] = v
	}

	// merge given headers with client-wide headers
	for k, v := range self.headers {
		finalHeaders[k] = v
	}

	for k, v := range headers {
		finalHeaders[k] = v
	}

	var reqUrl *url.URL
	var err error

	switch s, _ := stringutil.SplitPair(path, `://`); strings.ToLower(s) {
	case `http`, `https`:
		reqUrl, err = url.Parse(path)
	default:
		reqUrl, err = UrlPathJoin(self.uri, path)
	}

	if err == nil {
		// set querystring values from map
		for k, v := range finalParams {
			SetQ(reqUrl, k, v)
		}

		var payload io.Reader
		var encoded any

		if body != nil {
			if lit, ok := body.(Literal); ok {
				payload = bytes.NewReader([]byte(lit))
			} else if enc, err := self.encoder(body); err == nil {
				encoded = enc
				payload = enc
			} else {
				return nil, fmt.Errorf("encoder error: %v", err)
			}
		}

		if request, err := http.NewRequestWithContext(
			ctx,
			strings.ToUpper(string(method)),
			reqUrl.String(),
			payload,
		); err == nil {
			// Okay, this is a little weird...
			//
			// the existing signature for EncoderFunc accepts an interface with the original
			// intention of just having a stateless transformation of one type of data into the
			// encoded version of that data.
			//
			// However, it became clear that it would be useful to ALSO allow encoders to modify
			// the request data (e.g.: add headers, set Content-Type, etc.)
			//
			// So we're going to call encoder _again_, passing it the pointer to the newly-created
			// request.  We won't check the response for errors, and just assume that if an
			// EncoderFunc implementation supports doing something with *http.Request, it will
			// Do What It Needs To Do (TM) and move on.
			//
			// This is all in service of not breaking backwards compatibility by changing this
			// function signature.  Yay.
			//
			self.encoder(request)

			// finally, because this wasn't designed right from the start and should be using
			// a context.Context, we need to special case some things we know about
			if mpfr, ok := encoded.(*multipartFormRequest); ok {
				// special case: Multipart Form requests with the boundary value in the content type
				request.Header.Set(`Content-Type`, mpfr.contentType)
			}

			for k, v := range finalHeaders {
				request.Header.Set(k, typeutil.String(v))
			}

			// if we're going to auto-inject credentials from a .netrc file, attempt to do that now
			if self.basicAuthViaNetrc {
				if u, p, ok := NetrcCredentials(request.URL.Hostname()); ok {
					request.Header.Set(`Authorization`, EncodeBasicAuth(u, p))
				}
			}

			var hookObject any

			if self.preRequestHook != nil {
				if v, err := self.preRequestHook(request); err == nil {
					hookObject = v
				} else {
					return nil, fmt.Errorf("pre-request hook: %v", err)
				}
			}

			// close connection after sending this request and reading its response
			request.Close = true

			// perform the request
			if response, err := self.httpClient.Do(request); err == nil {
				if self.postRequestHook != nil {
					if err := self.postRequestHook(response, hookObject); err != nil {
						return nil, fmt.Errorf("post-request hook: %v", err)
					}
				}

				if response.StatusCode < 400 {
					return response, nil
				} else if self.errorDecoder != nil {
					if err := self.errorDecoder(response); err != nil {
						return response, err
					}
				}

				return response, fmt.Errorf("HTTP %v", response.Status)
			} else {
				return nil, fmt.Errorf("http request: %v", err)
			}
		} else {
			return nil, fmt.Errorf("request init error: %v", err)
		}
	} else {
		return nil, fmt.Errorf("url error: %v", err)
	}
}

func (self *Client) Encode(in any) ([]byte, error) {
	if self.encoder != nil {
		if r, err := self.encoder(in); err == nil {
			return io.ReadAll(r)
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("No encoder set")
	}
}

// Decode a response and, if applicable, automatically close the reader.
func (self *Client) Decode(r io.Reader, out any) error {
	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	if self.decoder != nil {
		return self.decoder(r, out)
	} else {
		return fmt.Errorf("No decoder set")
	}
}

func (self *Client) Get(path string, params map[string]any, headers map[string]any) (*http.Response, error) {
	return self.Request(Get, path, nil, params, headers)
}

func (self *Client) GetWithBody(path string, body any, params map[string]any, headers map[string]any) (*http.Response, error) {
	return self.Request(Get, path, body, params, headers)
}

func (self *Client) Post(path string, body any, params map[string]any, headers map[string]any) (*http.Response, error) {
	return self.Request(Post, path, body, params, headers)
}

func (self *Client) Put(path string, body any, params map[string]any, headers map[string]any) (*http.Response, error) {
	return self.Request(Put, path, body, params, headers)
}

func (self *Client) Delete(path string, params map[string]any, headers map[string]any) (*http.Response, error) {
	return self.Request(Delete, path, nil, params, headers)
}
