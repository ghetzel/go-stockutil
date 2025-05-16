package log

import (
	"errors"
	"testing"

	"github.com/ghetzel/testify/require"
)

type traceTest struct{}

func (self *traceTest) do(t *testing.T) {
	var trace = StackTrace(0)

	require.True(t, len(trace) >= 3)
	require.Equal(t, trace[0].Function, `runtime.Callers`)
	require.Equal(t, trace[0].FunctionName, `Callers`)
	require.Equal(t, trace[0].PackageName, `runtime`)
	require.Equal(t, trace[0].Receiver, ``)

	require.Equal(t, trace[1].Function, `github.com/ghetzel/go-stockutil/log.StackTrace`)
	require.Equal(t, trace[1].FunctionName, `StackTrace`)
	require.Equal(t, trace[1].PackageName, `github.com/ghetzel/go-stockutil/log`)
	require.Equal(t, trace[1].Receiver, ``)

	require.Equal(t, trace[2].Function, `github.com/ghetzel/go-stockutil/log.(*traceTest).do`)
	require.Equal(t, trace[2].FunctionName, `do`)
	require.Equal(t, trace[2].PackageName, `github.com/ghetzel/go-stockutil/log`)
	require.Equal(t, trace[2].Receiver, `(*traceTest)`)
}

func TestStackTrace(t *testing.T) {
	new(traceTest).do(t)
}

func TestErrors(t *testing.T) {
	var assert = require.New(t)

	var e1m = `error 1`
	var e1 = errors.New(e1m)

	var e2m = `error`
	var e2 = errors.New(e2m)

	assert.False(ErrContains(nil, nil))
	assert.False(ErrContains(e1, nil))
	assert.False(ErrContains(nil, e2))
	assert.False(ErrContains(nil, e2m))
	assert.True(ErrContains(e1, e1))
	assert.True(ErrContains(e1, e1m))
	assert.True(ErrContains(e2, e2))
	assert.True(ErrContains(e2, e2m))
	assert.True(ErrContains(e1, e2))
	assert.True(ErrContains(e1, e2m))
	assert.False(ErrContains(e2, e1))
	assert.False(ErrContains(e2, e1m))

	assert.True(ErrHasPrefix(e1, e1))
	assert.True(ErrHasPrefix(e1, e1m))
	assert.True(ErrHasPrefix(e2, e2))
	assert.True(ErrHasPrefix(e2, e2m))
	assert.True(ErrHasPrefix(e1, e2))
	assert.True(ErrHasPrefix(e1, e2m))
	assert.False(ErrHasPrefix(e1, `nope`))
	assert.False(ErrHasPrefix(e2, e1))
	assert.False(ErrHasPrefix(e2, e1m))

	assert.False(ErrHasSuffix(e1, e2))
	assert.False(ErrHasSuffix(e1, e2m))
	assert.True(ErrHasSuffix(e1, `1`))
	assert.False(ErrHasSuffix(e2, e1))
	assert.False(ErrHasSuffix(e2, e1m))
	assert.True(ErrHasSuffix(e1, e1))
	assert.True(ErrHasSuffix(e1, e1m))
	assert.True(ErrHasSuffix(e2, e2))
	assert.True(ErrHasSuffix(e2, e2m))
}
