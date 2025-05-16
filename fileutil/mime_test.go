package fileutil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestGetMimeType(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(`text/plain; charset=utf-8`, GetMimeType(`test.txt`))
	assert.Equal(`image/jpeg`, GetMimeType(`test.jpg`))
	assert.Equal(`image/gif`, GetMimeType(`test.gif`))
	assert.Equal(`image/png`, GetMimeType(`test.png`))
	assert.Equal(`image/svg+xml`, GetMimeType(`test.svg`))
	assert.Equal(`application/json`, GetMimeType(`test.json`))
	assert.Equal(`text/html; charset=utf-8`, GetMimeType(`test.html`))
	assert.Equal(`text/html; charset=utf-8`, GetMimeType(`test.htm`))
}
