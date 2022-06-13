package httputil

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/testify/assert"
	"github.com/ghetzel/testify/require"
)

func testHttpServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case `GET`:
			RespondJSON(w, map[string]interface{}{
				`path`: req.URL.Path,
				`qs`:   req.URL.Query(),
			})
		case `POST`:
			switch ct, _ := stringutil.SplitPair(req.Header.Get(`Content-Type`), `;`); ct {
			case `multipart/form-data`:
				if err := req.ParseMultipartForm(1048576); err == nil {
					RespondJSON(w, req.PostForm, http.StatusAccepted)
				} else {
					RespondJSON(w, err, http.StatusBadRequest)
				}

			default:
				var input interface{}

				if !QBool(req, `thing`) {
					RespondJSON(w, nil, http.StatusForbidden)
					return
				}

				if err := ParseJSON(req.Body, &input); err == nil {
					RespondJSON(w, input, http.StatusCreated)
				} else {
					RespondJSON(w, err, http.StatusBadRequest)
				}
			}

		case `PUT`:
			var input interface{}

			if !QBool(req, `thing`) || !QBool(req, `topthing`) {
				RespondJSON(w, nil, http.StatusForbidden)
				return
			}

			if err := ParseJSON(req.Body, &input); err == nil {
				RespondJSON(w, input, http.StatusAccepted)
			} else {
				RespondJSON(w, err, http.StatusBadRequest)
			}

		case `DELETE`:
			RespondJSON(w, nil)
		case `PATCH`:
			var err error

			if u, p, ok := req.BasicAuth(); ok {
				if u == `test` {
					if p == `test` {
						RespondJSON(w, nil, http.StatusAccepted)
						return
					} else {
						err = errors.New(`bad password`)
					}
				} else {
					err = errors.New(`bad username`)
				}
			} else {
				err = errors.New(`no auth provided`)
			}

			RespondJSON(w, err, http.StatusUnauthorized)
		}
	}))
}

func TestDefaultClient(t *testing.T) {
	assert := require.New(t)
	server := testHttpServer()
	defer server.Close()

	// GET
	// --------------------------------------------------------------------------------------------
	var data, err = GetBody(server.URL + `/test/path`)
	assert.NoError(err)
	assert.Equal([]byte("{\"path\":\"/test/path\",\"qs\":{}}\n"), data)
}

func TestClient(t *testing.T) {
	assert := require.New(t)
	var out map[string]interface{}
	var outS string

	server := testHttpServer()
	defer server.Close()

	client, err := NewClient(server.URL + `/base/?hello=true`)
	client.SetParam(`topthing`, true)

	assert.NoError(err)
	assert.NotNil(client)

	// GET
	// --------------------------------------------------------------------------------------------
	response, err := client.Get(`/test/path`, nil, nil)
	assert.NoError(err)
	assert.NotNil(response)
	assert.NoError(ParseJSON(response.Body, &out))
	assert.Equal(map[string]interface{}{
		`path`: `/base/test/path`,
		`qs`: map[string]interface{}{
			`hello`:    []interface{}{`true`},
			`topthing`: []interface{}{`true`},
		},
	}, out)

	// POST
	// --------------------------------------------------------------------------------------------
	response, err = client.Post(`/base/test/path`, `postable`, map[string]interface{}{
		`thing`: true,
	}, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.NoError(client.Decode(response.Body, &outS))
	assert.Equal(`postable`, outS)

	// PUT
	// --------------------------------------------------------------------------------------------
	response, err = client.Put(`/base/test/path`, `puttable`, map[string]interface{}{
		`thing`: true,
	}, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.NoError(client.Decode(response.Body, &outS))
	assert.Equal(`puttable`, outS)

	// PUT
	// --------------------------------------------------------------------------------------------
	response, err = client.Delete(`/base/test/path`, nil, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.Equal(http.StatusNoContent, response.StatusCode)
}

func TestClientMultipartFormEncoder(t *testing.T) {
	assert := require.New(t)
	var out map[string]interface{}

	server := testHttpServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	assert.NoError(err)
	assert.NotNil(client)

	client.SetErrorDecoder(func(res *http.Response) error {
		assert.NotNil(res.Body)

		return errors.New(typeutil.String(res.Body))
	})

	response, err := client.WithEncoder(MultipartFormEncoder).Post(`/way/cool`, map[string]interface{}{
		`file`:   bytes.NewBuffer([]byte("test file 1\n")),
		`other`:  bytes.NewBuffer([]byte("test file 2\n")),
		`key`:    `value`,
		`enable`: true,
	}, nil, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.Equal(response.StatusCode, http.StatusAccepted)
	assert.NoError(ParseJSON(response.Body, &out))

	assert.Equal(map[string]interface{}{
		`file`:   []interface{}{"test file 1\n"},
		`other`:  []interface{}{"test file 2\n"},
		`key`:    []interface{}{"value"},
		`enable`: []interface{}{"true"},
	}, out)
}

func TestClientNetrcAuth(t *testing.T) {
	server := testHttpServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client.SetErrorDecoder(func(res *http.Response) error {
		assert.NotNil(t, res.Body)

		return errors.New(typeutil.String(res.Body))
	})

	NetrcFile = fileutil.MustWriteTempFile(
		"machine * login test password test\n",
		"test-ghetzel-go-stockutil-httputil",
	)

	defer os.Remove(NetrcFile)

	client.SetAutomaticLogin(true)

	response, err := client.Request(Patch, `/`, nil, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.StatusCode, http.StatusAccepted)
}
