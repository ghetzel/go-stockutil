package fileutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/testify/assert"
)

func testHttpServer(t *testing.T, mustHeaders ...map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if len(mustHeaders) > 0 {
			for k, v := range maputil.M(mustHeaders[0]).MapNative() {
				assert.Equal(t, v, req.Header.Get(k), fmt.Sprintf("bad header %v: %q != %q", k, req.Header.Get(k), v))
			}
		}

		var _, act, arg = stringutil.SplitTriple(req.URL.Path, `/`)

		switch act {
		case `sleep`:
			time.Sleep(typeutil.Duration(arg))
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK"))
	}))
}

func TestRetrieve(t *testing.T) {
	var server = testHttpServer(t)
	defer server.Close()

	for uri, expected := range map[string]string{
		"testdir/a.txt":         "a\n",
		"./testdir/a.txt":       "a\n",
		"file:///testdir/a.txt": "a\n",
		server.URL + "/hello":   "OK",
	} {
		var rc, err = Retrieve(context.TODO(), uri)
		assert.NoError(t, err, uri)

		var actual, rerr = io.ReadAll(rc)

		assert.NoError(t, rc.Close(), uri)
		assert.NoError(t, rerr, uri)
		assert.Equal(t, expected, string(actual), uri)
	}
}
