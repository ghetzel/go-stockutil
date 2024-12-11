package typeutil

import (
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

func TestVariant(t *testing.T) {
	assert := require.New(t)

	assert.Equal(`test`, Variant{`test`}.String())
	assert.True(Variant{`True`}.Bool())
	assert.True(Variant{`true`}.Bool())
	assert.True(Variant{`TRUE`}.Bool())
	assert.True(Variant{`1`}.Bool())
	assert.True(Variant{`dennis`}.Bool())
	assert.True(Variant{`0.000000001`}.Bool())
	assert.False(Variant{`False`}.Bool())
	assert.False(Variant{`false`}.Bool())
	assert.False(Variant{`0`}.Bool())
	assert.False(Variant{`0.0`}.Bool())
	assert.Equal(int64(1), Variant{1}.Int())
	assert.Equal(int64(1), Variant{1.9}.Int())
	assert.Equal(float64(1.9), Variant{1.9}.Float())
	assert.True(time.Unix(1500000000, 0).Equal(Variant{1500000000}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`1500000000`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-14T02:40:00Z`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-14T02:40:00-00:00`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-13T22:40:00-04:00`}.Time()))
	assert.Equal([]byte{0x74, 0x65, 0x73, 0x74}, Variant{`test`}.Bytes())

	assert.Equal(map[Variant]Variant{
		V(`test`):  V(1),
		V(`other`): V(2.4),
	}, V(map[string]interface{}{
		`test`:  1,
		`other`: 2.4,
	}).Map())

	type vStructOne struct {
		Name    string
		Age     int
		Pi      float64
		enabled bool
	}

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(3.1415),
	}, V(vStructOne{
		Name:    `test`,
		Age:     42,
		Pi:      3.1415,
		enabled: true,
	}).Map())

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(3.1415),
	}, V(&vStructOne{
		Name:    `test`,
		Age:     42,
		Pi:      3.1415,
		enabled: true,
	}).Map())

	type vStructTagged struct {
		Name    string
		Age     int     `testaroo:"age"`
		Pi      float64 `testaroo:"pi,omitempty"`
		enabled bool
	}

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`age`):  V(42),
	}, V(vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map(`testaroo`))

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`age`):  V(42),
	}, V(&vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map(`testaroo`))

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(float64(0)),
	}, V(vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map())

	assert.Equal(map[Variant]Variant{
		V(`Name`): V(`test`),
		V(`Age`):  V(42),
		V(`Pi`):   V(float64(0)),
	}, V(&vStructTagged{
		Name:    `test`,
		Age:     42,
		enabled: true,
	}).Map())
}

func TestVariantAppend(t *testing.T) {
	var v = V(`hello`)

	require.NoError(t, v.Append(`there`))
	require.Equal(t, []string{`hello`, `there`}, v.Strings())
}

func TestNil(t *testing.T) {
	require.True(t, Nil().IsNil())
	require.True(t, Nil().IsZero())
}

func TestOrString(t *testing.T) {
	require.Equal(t, ``, OrString(``))
	require.Equal(t, `hey`, OrString(`hey`))
	require.Equal(t, `hey`, OrString(`hey`, ``, ``))
	require.Equal(t, `hey`, OrString(``, `hey`, ``))
	require.Equal(t, `hey`, OrString(``, ``, `hey`))
}

func TestOrBool(t *testing.T) {
	require.False(t, OrBool(false))
	require.False(t, OrBool(0))
	require.False(t, OrBool(`0`))
	require.False(t, OrBool(`false`))
	require.False(t, OrBool(`no`))
	require.False(t, OrBool(`off`))

	require.True(t, OrBool(true))
	require.True(t, OrBool(`true`))
	require.True(t, OrBool(1))
	require.True(t, OrBool(`1`))
	require.True(t, OrBool(`yes`))
	require.True(t, OrBool(`on`))
}

func TestOrInt(t *testing.T) {
	require.Equal(t, int64(42), OrInt(42, 96))
	require.Equal(t, int64(42), OrInt(`42`, 96))
	require.Equal(t, int64(42), OrInt(``, `0`, 42, 96))
	require.Equal(t, int64(42), OrInt(0, false, ``, 42, 96))
}

func TestOrFloat(t *testing.T) {
	require.Equal(t, float64(42), OrFloat(42, 96))
	require.Equal(t, float64(42), OrFloat(`42`, 96))
	require.Equal(t, float64(42), OrFloat(``, `0`, 42, 96))
	require.Equal(t, float64(42), OrFloat(0, false, ``, 42, 96))
}

func TestOrTime(t *testing.T) {
	require.True(t, OrTime(``).IsZero())
	require.True(t, OrTime(nil, ``, false, nil).IsZero())

	require.False(t, OrTime(`now`).IsZero())
	require.Equal(t, time.Unix(0, 0), OrTime(0))
}

func TestOrDuration(t *testing.T) {
	require.Zero(t, OrDuration(``))
	require.Zero(t, OrDuration(``, false, nil))

	require.Equal(t, 4*time.Hour, OrDuration(``, 0, false, `0ns`, `4h`))
	require.Equal(t, 24*time.Hour, OrDuration(`1d`))
	require.Equal(t, 5*time.Minute+3*time.Second, OrDuration(`5m3s`, `1m18s`))
}

func TestIsNumeric(t *testing.T) {
	require.True(t, IsNumeric(0))
	require.True(t, IsNumeric(`0`))
	require.True(t, IsNumeric(float64(0)))
	require.True(t, IsNumeric(1))
	require.True(t, IsNumeric(`1`))
	require.True(t, IsNumeric(float64(1)))

	require.False(t, IsNumeric(false))
	require.False(t, IsNumeric(true))
	require.False(t, IsNumeric(time.Now()))
	require.True(t, IsNumeric(time.Now().Unix()))
	require.False(t, IsNumeric(`a potato`))
	require.False(t, IsNumeric(``))
	require.False(t, IsNumeric(nil))
}

func TestIsTime(t *testing.T) {
	require.True(t, IsTime(0))
	require.True(t, IsTime(`0`))
	require.True(t, IsTime(float64(0)))
	require.True(t, IsTime(1))
	require.True(t, IsTime(`1`))
	require.True(t, IsTime(float64(1)))

	require.False(t, IsTime(false))
	require.False(t, IsTime(true))
	require.True(t, IsTime(time.Now()))
	require.True(t, IsTime(time.Now().Unix()))
	require.False(t, IsTime(`a potato`))
	require.True(t, IsTime(`now`))
	require.True(t, IsTime(`2006-01-02`))
	require.False(t, IsTime(``))
	require.False(t, IsTime(nil))
}

func TestIsDuration(t *testing.T) {
	require.False(t, IsDuration(0))
	require.False(t, IsDuration(`0`))
	require.False(t, IsDuration(float64(0)))
	require.False(t, IsDuration(1))
	require.False(t, IsDuration(`1`))
	require.False(t, IsDuration(float64(1)))

	require.False(t, IsDuration(false))
	require.False(t, IsDuration(true))
	require.False(t, IsDuration(`a potato`))
	require.True(t, IsDuration(`30s`))
	require.True(t, IsDuration(`20ms`))
	require.True(t, IsDuration(`5m`))
	require.False(t, IsDuration(`2006-01-02`))
	require.False(t, IsDuration(``))
	require.False(t, IsDuration(nil))
}

func TestGenericCompare(t *testing.T) {
	require.True(t, IsLessThan(false, true))
	require.True(t, IsLessThan(5, 10))
	require.True(t, IsLessThan(`apple`, `banana`))

	require.False(t, IsLessThan(nil, nil))
	require.False(t, IsLessThan(``, nil))
	require.False(t, IsLessThan(0, 0))
	require.False(t, IsLessThan(42, 42))
	require.False(t, IsLessThan(420, 69))
}

func TestVariantIsFunction(t *testing.T) {
	require.False(t, V(nil).IsFunction())
	require.False(t, V(false).IsFunction())
	require.False(t, V(nil).IsFunction(`func()`))
	require.False(t, V(false).IsFunction(`func()`))

	require.True(t, V(func() {}).IsFunction())
	require.True(t, V(func() {}).IsFunction(`func()`))
	require.True(t, V(func(_ string) {}).IsFunction(`func(string)`))
	require.True(t, V(func(_ string) error { return nil }).IsFunction(`func(string) error`))
	require.True(t, V(func(_ string) error { return nil }).IsFunction(`func(any) any`))
	require.True(t, V(func(_ string) (int, error) { return 0, nil }).IsFunction(`func(string) (int, error)`))
}

func TestSlices(t *testing.T) {
	require.Len(t, Slice(nil), 0)
	require.Len(t, Slice(1), 1)
	require.Len(t, Slice([]int{1, 2, 3}), 3)
}
