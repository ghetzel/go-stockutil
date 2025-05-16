package fileutil

import (
	"context"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/ghetzel/testify/assert"
)

func TestRetrieveViaHTTP(t *testing.T) {
	var server = testHttpServer(t, map[string]any{
		`X-Test`: `1`,
	})
	defer server.Close()

	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, `metadata`, map[string]string{
		`x-test`: `1`,
	})

	var u, err = url.Parse(server.URL)

	assert.NoError(t, err)

	var rc, rerr = RetrieveViaHTTP(ctx, u)
	assert.NoError(t, rerr)

	var data, derr = io.ReadAll(rc)

	assert.NoError(t, rc.Close())
	assert.NoError(t, derr)
	assert.Equal(t, "OK", string(data))

	u.Path = `/sleep/1100ms`

	_, rerr = RetrieveViaHTTP(ctx, u)
	assert.Contains(t, rerr.Error(), `context deadline exceeded`)
}
