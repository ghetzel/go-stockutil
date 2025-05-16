package utils

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestDetectConvertType(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(Nil, DetectConvertType(nil))
	assert.Equal(Nil, DetectConvertType(``))

	assert.Equal(String, DetectConvertType(`07753`))
	assert.Equal(String, DetectConvertType(`test`))

	assert.Equal(Time, DetectConvertType(`2015-05-01 00:15:16`))
	assert.Equal(Time, DetectConvertType(`Fri May 1 00:15:16 UTC 2015`))
	assert.Equal(Time, DetectConvertType(`Fri May 01 00:15:16 +0000 2015`))
	assert.Equal(Time, DetectConvertType(`01 May 15 00:15 UTC`))
	assert.Equal(Time, DetectConvertType(`01 May 15 00:15 +0000`))
	assert.Equal(Time, DetectConvertType(`Friday, 01-May-15 00:15:16 UTC`))
	assert.Equal(Time, DetectConvertType(`2003-06-08T11:56`))
	assert.Equal(Time, DetectConvertType(`2003-06-08T11:56:36`))
	assert.Equal(Time, DetectConvertType(`2003-06-08 11:56`))
	assert.Equal(Time, DetectConvertType(`2003-06-08 11:56:36`))

	assert.Equal(Integer, DetectConvertType(`1136239445`))
	assert.Equal(Integer, DetectConvertType(`1136239445000`))
	assert.Equal(Integer, DetectConvertType(`1136239445000000`))
	assert.Equal(Integer, DetectConvertType(`0`))
	assert.Equal(Integer, DetectConvertType(`1`))
	assert.Equal(Integer, DetectConvertType(`17753`))
	assert.Equal(Integer, DetectConvertType(`0xdeadbeef`))
	assert.Equal(Integer, DetectConvertType(`0xDEADBEEF`))
	assert.Equal(String, DetectConvertType(`deadbeef`))
	assert.Equal(String, DetectConvertType(`DEADBEEF`))

	assert.Equal(Float, DetectConvertType(`0.0`))
	assert.Equal(Float, DetectConvertType(`3.1415`))
	assert.Equal(Float, DetectConvertType(`3.0001`))
	assert.Equal(Float, DetectConvertType(`3.1000`))
}

func TestConvertToInteger(t *testing.T) {
	var assert = require.New(t)

	var i int64
	var err error

	i, err = ConvertToInteger(``)
	assert.NoError(err)
	assert.EqualValues(0, i)

	i, err = ConvertToInteger(`0`)
	assert.NoError(err)
	assert.EqualValues(0, i)

	i, err = ConvertToInteger(`123`)
	assert.NoError(err)
	assert.EqualValues(123, i)

	i, err = ConvertToInteger(`0x0`)
	assert.NoError(err)
	assert.EqualValues(0, i)

	i, err = ConvertToInteger(`0x1`)
	assert.NoError(err)
	assert.EqualValues(1, i)

	i, err = ConvertToInteger(`0xA`)
	assert.NoError(err)
	assert.EqualValues(10, i)

	i, err = ConvertToInteger(`0xF`)
	assert.NoError(err)
	assert.EqualValues(15, i)

	i, err = ConvertToInteger(`0x10`)
	assert.NoError(err)
	assert.EqualValues(16, i)

	i, err = ConvertToInteger(`0xG`)
	assert.NotNil(err)
}

func TestConvertTypeSpecificity(t *testing.T) {
	var assert = require.New(t)

	assert.False(Nil.IsSupersetOf(Nil))
	assert.False(Nil.IsSupersetOf(Bytes))
	assert.False(Nil.IsSupersetOf(String))
	assert.False(Nil.IsSupersetOf(Float))
	assert.False(Nil.IsSupersetOf(Integer))
	assert.False(Nil.IsSupersetOf(Time))
	assert.False(Nil.IsSupersetOf(Boolean))

	assert.False(Bytes.IsSupersetOf(Bytes))
	assert.True(Bytes.IsSupersetOf(String))
	assert.True(Bytes.IsSupersetOf(Float))
	assert.True(Bytes.IsSupersetOf(Integer))
	assert.True(Bytes.IsSupersetOf(Time))
	assert.True(Bytes.IsSupersetOf(Boolean))
	assert.True(Bytes.IsSupersetOf(Nil))

	assert.False(String.IsSupersetOf(Bytes))
	assert.False(String.IsSupersetOf(String))
	assert.True(String.IsSupersetOf(Float))
	assert.True(String.IsSupersetOf(Integer))
	assert.True(String.IsSupersetOf(Time))
	assert.True(String.IsSupersetOf(Boolean))
	assert.True(String.IsSupersetOf(Nil))

	assert.False(Float.IsSupersetOf(Bytes))
	assert.False(Float.IsSupersetOf(String))
	assert.False(Float.IsSupersetOf(Float))
	assert.True(Float.IsSupersetOf(Integer))
	assert.True(Float.IsSupersetOf(Time))
	assert.True(Float.IsSupersetOf(Boolean))
	assert.True(Float.IsSupersetOf(Nil))

	assert.False(Integer.IsSupersetOf(Bytes))
	assert.False(Integer.IsSupersetOf(String))
	assert.False(Integer.IsSupersetOf(Float))
	assert.False(Integer.IsSupersetOf(Integer))
	assert.True(Integer.IsSupersetOf(Time))
	assert.True(Integer.IsSupersetOf(Boolean))
	assert.True(Integer.IsSupersetOf(Nil))
}
