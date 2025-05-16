package httputil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghetzel/testify/require"
)

func mt(ct string) *http.Request {
	var req = httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, ct)
	return req
}

func TestMediaType(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(``, MediaType(mt(``)))
	assert.Equal(`text/plain`, MediaType(mt(`text/plain`)))
	assert.Equal(`text/plain`, MediaType(mt(`text/plain; charset=utf-8`)))
}

func TestIsMediaType(t *testing.T) {
	var assert = require.New(t)

	var req = mt(`text/plain; charset=utf-8`)

	assert.True(IsMediaType(req, `text/plain`))
	assert.True(IsMediaType(req, `text/plain`, `text/html`))
	assert.True(IsMediaType(req, `text/`))
	assert.False(IsMediaType(req))
	assert.False(IsMediaType(req, `text/html`))
}

func ExampleIsMediaType_singleMediaType() {
	var req = httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `text/html`))
	// Output: true
}

func ExampleIsMediaType_multipleMediaTypes() {
	var req = httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `text/html`, `text/plain`))
	// Output: true
}

func ExampleIsMediaType_mediaTypePrefix() {
	var req = httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `text/`))
	// Output: true
}

func ExampleIsMediaType_nonMatchingPrefix() {
	var req = httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `video/`))
	// Output: false
}

func ExampleUrlPathJoin_simpleJoin() {
	url, _ := UrlPathJoin(`https://google.com`, `/search`)
	fmt.Println(url.String())
	// Output: https://google.com/search
}

func ExampleUrlPathJoin_slashJoin() {
	url, _ := UrlPathJoin(`https://google.com/`, `/`)
	fmt.Println(url.String())
	// Output: https://google.com/
}

func ExampleUrlPathJoin_emptyJoin() {
	url, _ := UrlPathJoin(`https://google.com/`, ``)
	fmt.Println(url.String())
	// Output: https://google.com/
}

func ExampleUrlPathJoin_joinWithQueryString() {
	url, _ := UrlPathJoin(`https://google.com/`, `/search?q=hello`)
	fmt.Println(url.String())
	// Output: https://google.com/search?q=hello
}

func ExampleUrlPathJoin_joinToExistingQueryStrings() {
	url, _ := UrlPathJoin(`https://google.com/search`, `?q=hello`)
	fmt.Println(url.String())
	// Output: https://google.com/search?q=hello
}

func ExampleUrlPathJoin_pathAndQuery() {
	url, _ := UrlPathJoin(`https://example.com/api/v1?hello=there`, `/things/new?example=true`)
	fmt.Println(url.String())
	// Output: https://example.com/api/v1/things/new?example=true&hello=there
}

func ExampleUrlPathJoin_pathAndQueryTrailingSlash() {
	url, _ := UrlPathJoin(`https://example.com/api/v1?hello=there`, `/things/new/?example=true`)
	fmt.Println(url.String())
	// Output: https://example.com/api/v1/things/new/?example=true&hello=there
}
