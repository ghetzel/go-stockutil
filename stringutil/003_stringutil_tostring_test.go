package stringutil

import (
	"fmt"
	"testing"

	"github.com/ghetzel/testify/require"
)

type fmtStringerTest struct{}

func (self *fmtStringerTest) String() string {
	return `stringer-worked`
}

func TestConvertToString(t *testing.T) {
	var assert = require.New(t)

	v, err := ConvertTo(String, nil)
	assert.NoError(err)
	assert.Equal(``, v)

	v, err = ConvertTo(String, []byte{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, v)

	v, err = ConvertTo(String, []uint8{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, v)
}
func TestToString(t *testing.T) {
	var assert = require.New(t)

	var testvalues = map[any]string{
		nil:                        ``,
		fmt.Errorf(`error-worked`): `error-worked`,
		new(fmtStringerTest):       `stringer-worked`,

		int(0):   `0`,
		int(1):   `1`,
		int8(0):  `0`,
		int8(1):  `1`,
		int16(0): `0`,
		int16(1): `1`,
		int32(0): `0`,
		int32(1): `1`,
		int64(0): `0`,
		int64(1): `1`,

		uint(0):   `0`,
		uint(1):   `1`,
		uint8(0):  `0`,
		uint8(1):  `1`,
		uint16(0): `0`,
		uint16(1): `1`,
		uint32(0): `0`,
		uint32(1): `1`,
		uint64(0): `0`,
		uint64(1): `1`,

		float32(0.0): `0`,
		float32(1.0): `1`,
		float64(0.0): `0`,
		float64(1.0): `1`,
		float32(0.5): `0.5`,
		float32(1.7): `1.7`,
		float64(0.6): `0.6`,
		float64(1.2): `1.2`,
	}

	for in, out := range testvalues {
		if v, err := ToString(in); err != nil || v != out {
			t.Errorf("Value %v (%T) ToString failed: expected '%s', got '%s' (err: %v)", in, in, out, v, err)
		}
	}

	s, err := ToString([]uint8{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, s)

	s, err = ToString([]byte{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, s)

	s, err = ToString([]rune{'t', 'e', 's', 't'})
	assert.NoError(err)
	assert.Equal(`test`, s)
}

func TestToStringSlice(t *testing.T) {
	var assert = require.New(t)

	v, err := ToStringSlice([]string{`a`, `b`, `c`})
	assert.Nil(err)
	assert.Equal([]string{`a`, `b`, `c`}, v)

	v, err = ToStringSlice([]int{1, 2, 3})
	assert.Nil(err)
	assert.Equal([]string{`1`, `2`, `3`}, v)

	v, err = ToStringSlice([]float32{1.5, 2.7, 3.0032477})
	assert.Nil(err)
	assert.Equal([]string{`1.5`, `2.7`, `3.0032477`}, v)

	v, err = ToStringSlice([]float64{1.5, 2.7, 3.9832412754892137})
	assert.Nil(err)
	assert.Equal([]string{`1.5`, `2.7`, `3.9832412754892137`}, v)

	v, err = ToStringSlice([]any{1, true, 3.9832412754892137, `four`})
	assert.Nil(err)
	assert.Equal([]string{`1`, `true`, `3.9832412754892137`, `four`}, v)

	v, err = ToStringSlice(true)
	assert.Nil(err)
	assert.Equal([]string{`true`}, v)

	v, err = ToStringSlice(true)
	assert.Nil(err)
	assert.Equal([]string{`true`}, v)

	v, err = ToStringSlice(nil)
	assert.Nil(err)
	assert.Equal([]string{}, v)
}
