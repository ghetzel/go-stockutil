package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsKind(t *testing.T) {
	var assert = require.New(t)

	assert.False(IsKind(nil, reflect.String))
	assert.False(IsKind(reflect.ValueOf(nil), reflect.String))
}
