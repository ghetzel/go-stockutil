package utils

import (
	"reflect"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestIsKind(t *testing.T) {
	var assert = require.New(t)

	assert.False(IsKind(nil, reflect.String))
	assert.False(IsKind(reflect.ValueOf(nil), reflect.String))
}
