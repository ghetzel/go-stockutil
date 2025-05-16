package typeutil

import (
	"errors"
	"strings"
	"testing"

	"github.com/ghetzel/testify/require"
)

type testCustomType struct{}

func TestFunctionArity(t *testing.T) {
	var assert = require.New(t)

	i, o, err := FunctionArity(strings.Compare)
	assert.NoError(err)
	assert.Equal(2, i)
	assert.Equal(1, o)

	var f0_0 = func() {}
	i, o, err = FunctionArity(f0_0)
	assert.NoError(err)
	assert.Equal(0, i)
	assert.Equal(0, o)

	var f0_1 = func() error { return nil }
	i, o, err = FunctionArity(f0_1)
	assert.NoError(err)
	assert.Equal(0, i)
	assert.Equal(1, o)
}

func TestParseSignatureString(t *testing.T) {
	//----------------------------------------------------------------------------------------------
	var ident, args, returns, err = ParseSignatureString(`testFunc(str,bool, *testCustomType) (bool,error)`)
	require.NoError(t, err)
	require.Equal(t, `testFunc`, ident)
	require.Len(t, args, 3)
	args[0].IsSameTypeAs(`example string`)
	args[1].IsSameTypeAs(true)
	args[2].IsSameTypeAs(new(testCustomType))
	require.Len(t, returns, 2)
	returns[0].IsSameTypeAs(true)
	returns[1].IsSameTypeAs(errors.New(`test error`))

	//----------------------------------------------------------------------------------------------
	ident, args, returns, err = ParseSignatureString(`func()`)
	require.NoError(t, err)
	require.Equal(t, `(anonymous)`, ident)
	require.Len(t, args, 0)
	require.Len(t, returns, 0)

	//----------------------------------------------------------------------------------------------
	ident, args, returns, err = ParseSignatureString(`testFunc() error`)
	require.NoError(t, err)
	require.Equal(t, `testFunc`, ident)
	require.Len(t, args, 0)
	require.Len(t, returns, 1)
	returns[0].IsSameTypeAs(errors.New(`test error`))

	//----------------------------------------------------------------------------------------------
	ident, args, returns, err = ParseSignatureString(`testFunc(bool) (error)`)
	require.NoError(t, err)
	require.Equal(t, `testFunc`, ident)
	require.Len(t, args, 1)
	args[0].IsSameTypeAs(false)
	require.Len(t, returns, 1)
	returns[0].IsSameTypeAs(errors.New(`test error`))
}

func TestFunctionMatchesSignature(t *testing.T) {
	require.NoError(t, FunctionMatchesSignature(func() {}, `func()`))
	require.NoError(t, FunctionMatchesSignature(func(_ string) {}, `func(string)`))
	require.NoError(t, FunctionMatchesSignature(func(_ string) error { return nil }, `func(string) error`))
	require.NoError(t, FunctionMatchesSignature(func(_ string) error { return nil }, `func(any) any`))
	require.NoError(t, FunctionMatchesSignature(func(_ string) (int, error) { return 0, nil }, `func(string) (int, error)`))
}
