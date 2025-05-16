package httputil

import (
	"net/http"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestAttachAndRetrieveValue(t *testing.T) {
	var assert = require.New(t)

	req, err := http.NewRequest(`GET`, `about:blank`, nil)
	assert.NoError(err)

	RequestSetValue(req, `test-value`, `123456789`)

	assert.False(RequestGetValue(req, `test-value`).IsNil())
	assert.False(RequestGetValue(req, `test-value`).IsZero())
	assert.EqualValues(123456789, RequestGetValue(req, `test-value`).Int())
	assert.EqualValues(`123456789`, RequestGetValue(req, `test-value`).String())
	assert.True(RequestGetValue(req, `test-value`).Bool())
}
