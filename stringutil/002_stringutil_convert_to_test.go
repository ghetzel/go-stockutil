package stringutil

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

func TestConvertToFloat(t *testing.T) {
	var assert = require.New(t)

	type testFloat float64
	const (
		testFloatZero testFloat = iota
		testFloatE              = 2.71828
		testFloatPi             = 3.141597
	)

	v, err := ConvertTo(Float, nil)
	assert.NoError(err)
	assert.Equal(float64(0), v)

	v, err = ConvertTo(Float, "1.5")
	assert.NoError(err)
	assert.Equal(float64(1.5), v)

	v, err = ConvertTo(Float, "1")
	assert.NoError(err)
	assert.Equal(float64(1.0), v)

	v, err = ConvertToFloat("1.5")
	assert.NoError(err)
	assert.Equal(float64(1.5), v)

	v, err = ConvertToFloat("1.0")
	assert.NoError(err)
	assert.Equal(float64(1.0), v)

	v, err = ConvertTo(Float, testFloatZero)
	assert.NoError(err)
	assert.Equal(float64(0), v)

	v, err = ConvertTo(Float, testFloatE)
	assert.NoError(err)
	assert.Equal(float64(2.71828), v)

	v, err = ConvertTo(Float, testFloatPi)
	assert.NoError(err)
	assert.Equal(float64(3.141597), v)

	for _, fail := range []string{`potato`, `true`, `2015-05-01 00:15:16`} {
		_, err := ConvertTo(Float, fail)
		assert.Error(err)

		_, err = ConvertToFloat(fail)
		assert.Error(err)
	}
}

func TestConvertToInteger(t *testing.T) {
	var assert = require.New(t)

	v, err := ConvertTo(Integer, nil)
	assert.NoError(err)
	assert.Equal(int64(0), v)

	v, err = ConvertTo(Integer, "7")
	assert.NoError(err)
	assert.Equal(int64(7), v)

	v, err = ConvertToInteger("7")
	assert.NoError(err)
	assert.Equal(int64(7), v)

	var tm = time.Date(2010, 2, 21, 15, 14, 13, 0, time.UTC)

	v, err = ConvertTo(Integer, tm)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	v, err = ConvertToInteger(tm)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	v, err = ConvertTo(Integer, `2010-02-21 15:14:13`)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	type testInt int64
	const (
		testInt1 testInt = iota
		testInt2         = 2
		testInt3         = 3
	)

	v, err = ConvertTo(Integer, testInt1)
	assert.NoError(err)
	assert.Equal(int64(0), v)

	v, err = ConvertTo(Integer, testInt2)
	assert.NoError(err)
	assert.Equal(int64(2), v)

	v, err = ConvertTo(Integer, testInt3)
	assert.NoError(err)
	assert.Equal(int64(3), v)

	for _, fail := range []string{`0.0`, `1.5`, `potato`, `true`} {
		_, err := ConvertTo(Integer, fail)
		assert.Error(err)

		_, err = ConvertToInteger(fail)
		assert.Error(err)
	}
}

func TestConvertToBytes(t *testing.T) {
	var assert = require.New(t)

	v, err := ConvertTo(Bytes, nil)
	assert.NoError(err)
	assert.Equal([]byte{}, v)

	v, err = ConvertTo(Bytes, []byte{})
	assert.NoError(err)
	assert.Equal([]byte{}, v)

	v, err = ConvertTo(Bytes, []byte{1, 2, 3})
	assert.NoError(err)
	assert.Equal([]byte{1, 2, 3}, v)

	v, err = ConvertTo(Bytes, `test`)
	assert.NoError(err)
	assert.Equal([]byte{0x74, 0x65, 0x73, 0x74}, v)

	v, err = ConvertTo(Bytes, []int{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, string(v.([]byte)))
}

func TestConvertToBoolean(t *testing.T) {
	var assert = require.New(t)

	v, err := ConvertTo(Boolean, nil)
	assert.Equal(false, v)

	v, err = ConvertTo(Boolean, `true`)
	assert.NoError(err)
	assert.Equal(true, v)

	v, err = ConvertTo(Boolean, `false`)
	assert.NoError(err)
	assert.Equal(false, v)

	v, err = ConvertToBool(`true`)
	assert.NoError(err)
	assert.Equal(true, v)

	v, err = ConvertToBool(`false`)
	assert.NoError(err)
	assert.Equal(false, v)

	for _, fail := range []string{`1.5`, `potato`, `01`, `2015-05-01 00:15:16`} {
		_, err := ConvertTo(Boolean, fail)
		assert.Error(err)

		_, err = ConvertToBool(fail)
		assert.Error(err)
	}
}

func TestConvertToTime(t *testing.T) {
	var assert = require.New(t)

	var atLeastNow = time.Now()

	var values = map[string]time.Time{
		`2015-05-01 00:15:16`:         time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		`Fri May 1 00:15:16 UTC 2015`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `Fri May 01 00:15:16 +0000 2015`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `01 May 15 00:15 UTC`:            time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `01 May 15 00:15 +0000`:          time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `Friday, 01-May-15 00:15:16 UTC`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		`1136239445`:          time.Date(2006, 1, 2, 17, 4, 5, 0, time.Now().Location()),
		`2038-01-19 03:14:06`: time.Date(2038, 1, 19, 3, 14, 6, 0, time.UTC),
		`2038-01-19 03:14:07`: time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC),
		`2038-01-19 03:14:08`: time.Date(2038, 1, 19, 3, 14, 8, 0, time.UTC),
		`2147483646`:          time.Date(2038, 1, 19, 3, 14, 6, 0, time.UTC),
		`2147483647`:          time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC),
		`2147483648`:          time.Date(2038, 1, 19, 3, 14, 8, 0, time.UTC),
	}

	v, err := ConvertToTime(`now`)
	assert.Nil(err)
	assert.True(v.After(atLeastNow))

	v, err = ConvertToTime(time.Now())
	assert.Nil(err)
	assert.True(v.After(atLeastNow))

	v, err = ConvertToTime(`0000-00-00 00:00:00`)
	assert.Nil(err)
	assert.Zero(v)

	for in, out := range values {
		v, err := ConvertTo(Time, in)
		assert.NoError(err)
		assert.IsType(time.Now(), v)
		assert.True(out.Equal(v.(time.Time)), in)

		v, err = ConvertToTime(in)
		assert.NoError(err)
		assert.True(out.Equal(v.(time.Time)), in)
	}

	for _, fail := range []string{`1.5`, `potato`, `false`} {
		_, err := ConvertTo(Time, fail)
		assert.Error(err)

		_, err = ConvertToTime(fail)
		assert.Error(err)
	}
}

func TestAutotypeNil(t *testing.T) {
	var assert = require.New(t)

	for _, testValue := range []string{
		``,
		`nil`,
		`null`,
		`Nil`,
		`NULL`,
		`None`,
		`undefined`,
	} {
		assert.Nil(Autotype(testValue), fmt.Sprintf("%q was not autotyped to nil", testValue))
	}
}

func TestAutotypeFloat(t *testing.T) {
	var assert = require.New(t)

	for _, testValue := range []string{
		`-0.00000000001`,
		`0.00000000001`,
		`1.5`,
		`-1.5`,
		fmt.Sprintf("%f", math.MaxFloat64),
		fmt.Sprintf("%f", -1*math.MaxFloat64),
	} {
		assert.IsType(float64(0), Autotype(testValue), testValue)
	}
}

func TestAutotypeInt(t *testing.T) {
	var assert = require.New(t)

	for _, testValue := range []string{
		`-1`,
		`0`,
		`1`,
		`12345`,
		`-12345`,
		fmt.Sprintf("%d", math.MaxInt64),
		fmt.Sprintf("%d", math.MinInt64),
	} {
		assert.IsType(int64(0), Autotype(testValue))
	}
}

func TestAutotypePreserveLeadingZeroes(t *testing.T) {
	var assert = require.New(t)

	for _, testValue := range []string{
		`00`,
		`01`,
		`07753`,
		`06094`,
		`0000000010000000`,
	} {
		assert.IsType(``, Autotype(testValue))
	}
}

func TestAutotypeDate(t *testing.T) {
	var assert = require.New(t)

	for _, testValue := range TimeFormats {
		var tvS = strings.Replace(string(testValue), `_`, ``, -1)
		tvS = strings.TrimSuffix(tvS, `07:00`)
		assert.IsType(time.Now(), Autotype(tvS))
	}
}

func TestAutotypeBool(t *testing.T) {
	var assert = require.New(t)

	for _, testValue := range []string{
		`true`,
		`True`,
		`false`,
		`False`,
	} {
		assert.IsType(true, Autotype(testValue))
	}

	for _, testValue := range []string{
		`trues`,
		`Falses`,
		`potato`,
	} {
		assert.IsType(``, Autotype(testValue))
	}
}
