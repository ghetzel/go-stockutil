package colorutil

import (
	"fmt"
	"testing"

	"github.com/ghetzel/testify/assert"
	"github.com/ghetzel/testify/require"
)

func assertHslToRgb(t *testing.T, h, s, l, r, g, b float64) {
	var r1, g1, b1 float64

	var assert = assert.New(t)

	r1, g1, b1 = hsl2rgb(h, s, l)

	assert.Equal(r, r1, fmt.Sprintf("hsl(%v,%v,%v) -> red(%v)", h, s, l, r1))
	assert.Equal(g, g1, fmt.Sprintf("hsl(%v,%v,%v) -> green(%v)", h, s, l, g1))
	assert.Equal(b, b1, fmt.Sprintf("hsl(%v,%v,%v) -> blue(%v)", h, s, l, b1))
}

func assertHsvToRgb(t *testing.T, h, s, v, r, g, b float64) {
	var r1, g1, b1 float64

	var assert = assert.New(t)

	r1, g1, b1 = hsv2rgb(h, s, v)

	assert.Equal(r, r1, fmt.Sprintf("hsv(%v,%v,%v) -> red(%v)", h, s, v, r1))
	assert.Equal(g, g1, fmt.Sprintf("hsv(%v,%v,%v) -> green(%v)", h, s, v, g1))
	assert.Equal(b, b1, fmt.Sprintf("hsv(%v,%v,%v) -> blue(%v)", h, s, v, b1))
}

func assertRgbToHsl(t *testing.T, r, g, b, h, s, l float64) {
	var h1, s1, l1 float64

	var assert = assert.New(t)

	h1, s1, l1 = rgb2lhs(HSL, r, g, b)

	assert.Equal(h, h1, fmt.Sprintf("rgb(%v,%v,%v) -> hue(%v)", r, g, b, h1))
	assert.Equal(s, s1, fmt.Sprintf("rgb(%v,%v,%v) -> saturation(%v)", r, g, b, s1))
	assert.Equal(l, l1, fmt.Sprintf("rgb(%v,%v,%v) -> lightness(%v)", r, g, b, l1))
}

func assertColor(t *testing.T, in string, rIn uint8, gIn uint8, bIn uint8, aIn uint8) {
	var assert = assert.New(t)

	var color Color
	var err error
	var r, g, b, a uint8

	color, err = Parse(in)
	assert.NoError(err)

	r, g, b, a = color.RGBA255()
	assert.Equal(rIn, r, fmt.Sprintf("%v: red", in))
	assert.Equal(gIn, g, fmt.Sprintf("%v: green", in))
	assert.Equal(bIn, b, fmt.Sprintf("%v: blue", in))
	assert.Equal(aIn, a, fmt.Sprintf("%v: alpha", in))
}

func adjustHue(t *testing.T, in any, adjustAmount float64, wantedColor string) {
	var assert = assert.New(t)

	var color Color
	var err error

	color, err = AdjustHue(in, adjustAmount)
	assert.NoError(err)

	assert.True(color.Equals(wantedColor), fmt.Sprintf("have: %v, want: %v", color, wantedColor))
}

func TestHslToRgb(t *testing.T) {
	// greyscale tests
	assertHslToRgb(t, 0, 0, 0, 0, 0, 0)
	assertHslToRgb(t, 0, 0, 0.1, 0.1, 0.1, 0.1)
	assertHslToRgb(t, 0, 0, 0.5, 0.5, 0.5, 0.5)
	assertHslToRgb(t, 0, 0, 0.75, 0.75, 0.75, 0.75)
	assertHslToRgb(t, 0, 0, 1, 1, 1, 1)

	// redsies
	assertHslToRgb(t, 0, 1, 0.01, 0.02, 0, 0)
	assertHslToRgb(t, 0, 1, 0.02, 0.04, 0, 0)
	assertHslToRgb(t, 0, 1, 0.04, 0.08, 0, 0)
	assertHslToRgb(t, 0, 1, 0.08, 0.16, 0, 0)
	assertHslToRgb(t, 0, 1, 0.16, 0.32, 0, 0)
	assertHslToRgb(t, 0, 1, 0.17, 0.34, 0, 0)
	assertHslToRgb(t, 0, 1, 0.31, 0.62, 0, 0)
	assertHslToRgb(t, 0, 1, 0.5, 1, 0, 0)

	// greensies
	assertHslToRgb(t, 120, 1, 0.01, 0, 0.02, 0)
	assertHslToRgb(t, 120, 1, 0.02, 0, 0.04, 0)
	assertHslToRgb(t, 120, 1, 0.04, 0, 0.08, 0)
	assertHslToRgb(t, 120, 1, 0.08, 0, 0.16, 0)
	assertHslToRgb(t, 120, 1, 0.16, 0, 0.32, 0)
	assertHslToRgb(t, 120, 1, 0.17, 0, 0.34, 0)
	assertHslToRgb(t, 120, 1, 0.31, 0, 0.62, 0)
	assertHslToRgb(t, 120, 1, 0.5, 0, 1, 0)

	// bluesies
	assertHslToRgb(t, 240, 1, 0.01, 0, 0, 0.02)
	assertHslToRgb(t, 240, 1, 0.02, 0, 0, 0.04)
	assertHslToRgb(t, 240, 1, 0.04, 0, 0, 0.08)
	assertHslToRgb(t, 240, 1, 0.08, 0, 0, 0.16)
	assertHslToRgb(t, 240, 1, 0.16, 0, 0, 0.32)
	assertHslToRgb(t, 240, 1, 0.17, 0, 0, 0.34)
	assertHslToRgb(t, 240, 1, 0.31, 0, 0, 0.62)
	assertHslToRgb(t, 240, 1, 0.5, 0, 0, 1)
}

func TestRgbToHsl(t *testing.T) {
	// greyscale tests
	assertRgbToHsl(t, 0, 0, 0, 0, 0, 0)
	assertRgbToHsl(t, 0.1, 0.1, 0.1, 0, 0, 0.1)
	assertRgbToHsl(t, 0.5, 0.5, 0.5, 0, 0, 0.5)
	assertRgbToHsl(t, 0.75, 0.75, 0.75, 0, 0, 0.75)
	assertRgbToHsl(t, 1, 1, 1, 0, 0, 1)

	// redsies
	assertRgbToHsl(t, 0.02, 0, 0, 0, 1, 0.01)
	assertRgbToHsl(t, 0.04, 0, 0, 0, 1, 0.02)
	assertRgbToHsl(t, 0.08, 0, 0, 0, 1, 0.04)
	assertRgbToHsl(t, 0.16, 0, 0, 0, 1, 0.08)
	assertRgbToHsl(t, 0.32, 0, 0, 0, 1, 0.16)
	assertRgbToHsl(t, 0.34, 0, 0, 0, 1, 0.17)
	assertRgbToHsl(t, 0.62, 0, 0, 0, 1, 0.31)
	assertRgbToHsl(t, 1, 0, 0, 0, 1, 0.5)

	// greensies
	assertRgbToHsl(t, 0, 0.02, 0, 120, 1, 0.01)
	assertRgbToHsl(t, 0, 0.04, 0, 120, 1, 0.02)
	assertRgbToHsl(t, 0, 0.08, 0, 120, 1, 0.04)
	assertRgbToHsl(t, 0, 0.16, 0, 120, 1, 0.08)
	assertRgbToHsl(t, 0, 0.32, 0, 120, 1, 0.16)
	assertRgbToHsl(t, 0, 0.34, 0, 120, 1, 0.17)
	assertRgbToHsl(t, 0, 0.62, 0, 120, 1, 0.31)
	assertRgbToHsl(t, 0, 1, 0, 120, 1, 0.5)

	// bluesies
	assertRgbToHsl(t, 0, 0, 0.02, 240, 1, 0.01)
	assertRgbToHsl(t, 0, 0, 0.04, 240, 1, 0.02)
	assertRgbToHsl(t, 0, 0, 0.08, 240, 1, 0.04)
	assertRgbToHsl(t, 0, 0, 0.16, 240, 1, 0.08)
	assertRgbToHsl(t, 0, 0, 0.32, 240, 1, 0.16)
	assertRgbToHsl(t, 0, 0, 0.34, 240, 1, 0.17)
	assertRgbToHsl(t, 0, 0, 0.62, 240, 1, 0.31)
	assertRgbToHsl(t, 0, 0, 1, 240, 1, 0.5)
}

func TestHsvToRgbFloat(t *testing.T) {
	// greyscale tests
	assertHsvToRgb(t, 0, 0, 0, 0, 0, 0)
	assertHsvToRgb(t, 0, 0, 0.1, 0.1, 0.1, 0.1)
	assertHsvToRgb(t, 0, 0, 0.5, 0.5, 0.5, 0.5)
	assertHsvToRgb(t, 0, 0, 0.75, 0.75, 0.75, 0.75)
	assertHsvToRgb(t, 0, 0, 1, 1, 1, 1)

	// redsies
	assertHsvToRgb(t, 0, 1, 1, 1, 0, 0)
	assertHsvToRgb(t, 0, 1, 0.75, 0.75, 0, 0)
	assertHsvToRgb(t, 0, 1, 0.5, 0.5, 0, 0)
	assertHsvToRgb(t, 0, 1, 0.01, 0.01, 0, 0)

	// greensies
	assertHsvToRgb(t, 120, 1, 1, 0, 1, 0)
	assertHsvToRgb(t, 120, 1, 0.75, 0, 0.75, 0)
	assertHsvToRgb(t, 120, 1, 0.5, 0, 0.5, 0)
	assertHsvToRgb(t, 120, 1, 0.01, 0, 0.01, 0)

	// bluesies
	assertHsvToRgb(t, 240, 1, 1, 0, 0, 1)
	assertHsvToRgb(t, 240, 1, 0.75, 0, 0, 0.75)
	assertHsvToRgb(t, 240, 1, 0.5, 0, 0, 0.5)
	assertHsvToRgb(t, 240, 1, 0.01, 0, 0, 0.01)
}

func TestRgbToHslSymmetry(t *testing.T) {
	var assert = require.New(t)
	var h, s, l float64
	var r, g, b uint8

	h, s, l = RgbToHsl(255, 0, 204)
	assert.Equal(float64(312), h)
	assert.Equal(float64(1), s)
	assert.Equal(float64(0.5), l)

	// verify that the operation is symmetric
	r, g, b = HslToRgb(h, s, l)
	assert.Equal(uint8(255), r)
	assert.Equal(uint8(0), g)
	assert.Equal(uint8(204), b)
}

func TestParse(t *testing.T) {
	assertColor(t, `#FFFFFF`, 255, 255, 255, 255)
	assertColor(t, `#FFFF00`, 255, 255, 0, 255)
	assertColor(t, `#FF00CC`, 255, 0, 204, 255)
	assertColor(t, `#FFFFFF00`, 255, 255, 255, 0)
	assertColor(t, `#FFFF00FF`, 255, 255, 0, 255)
	assertColor(t, `#FF00FFFF`, 255, 0, 255, 255)
	assertColor(t, `#00FFFFFF`, 0, 255, 255, 255)
	assertColor(t, `#FF000000`, 255, 0, 0, 0)
	assertColor(t, `#FF00FF00`, 255, 0, 255, 0)
	assertColor(t, `#00000000`, 0, 0, 0, 0)

	assertColor(t, `rgb(255, 255, 255)`, 255, 255, 255, 255)
	assertColor(t, `rgb(255, 255, 42)`, 255, 255, 42, 255)
	assertColor(t, `rgb(255, 42, 255)`, 255, 42, 255, 255)
	assertColor(t, `rgb(42, 255, 255)`, 42, 255, 255, 255)

	assertColor(t, `rgba(255, 255, 255, 255)`, 255, 255, 255, 255)
	assertColor(t, `rgba(255, 255, 255, 42)`, 255, 255, 255, 42)
	assertColor(t, `rgba(255, 255, 42, 255)`, 255, 255, 42, 255)
	assertColor(t, `rgba(255, 42, 255, 255)`, 255, 42, 255, 255)
	assertColor(t, `rgba(42, 255, 255, 255)`, 42, 255, 255, 255)

	assertColor(t, `rgba(255, 255, 255, 1.0)`, 255, 255, 255, 255)
	assertColor(t, `rgba(255, 255, 255, 0.5)`, 255, 255, 255, 128)
	assertColor(t, `rgba(255, 255, 255, 0.25)`, 255, 255, 255, 64)
	assertColor(t, `rgba(255, 255, 255, 0)`, 255, 255, 255, 0)

	assertColor(t, `hsl(0, 100%, 100%)`, 255, 255, 255, 255)
	assertColor(t, `hsl(0, 1, 1)`, 255, 255, 255, 255)
	assertColor(t, `hsl(0, 255, 255)`, 255, 255, 255, 255)
	assertColor(t, `hsla(0, 1, 1, 0.5)`, 255, 255, 255, 128)

	assertColor(t, `hsl(0, 100%, 50%)`, 255, 0, 0, 255)
	assertColor(t, `hsl(0, 100%, 25%)`, 128, 0, 0, 255)
	assertColor(t, `hsl(0, 100%, 12.5%)`, 64, 0, 0, 255)
	assertColor(t, `hsl(0, 100%, 6.25%)`, 33, 0, 0, 255)
	assertColor(t, `hsl(0, 100%, 3.125%)`, 15, 0, 0, 255)
	assertColor(t, `hsl(0, 100%, 1.5625%)`, 8, 0, 0, 255)

	assertColor(t, `hsl(0deg, 100%, 0.5)`, 255, 0, 0, 255)
	assertColor(t, `hsl(0deg, 100%, 0.25)`, 128, 0, 0, 255)
	assertColor(t, `hsl(0deg, 100%, 0.125)`, 64, 0, 0, 255)
	assertColor(t, `hsl(0deg, 100%, 0.0625)`, 33, 0, 0, 255)
	assertColor(t, `hsl(0deg, 100%, 0.03125)`, 15, 0, 0, 255)
	assertColor(t, `hsl(0deg, 100%, 0.015625)`, 8, 0, 0, 255)

	assertColor(t, `hsl(0, 255, 0.5)`, 255, 0, 0, 255)
	assertColor(t, `hsl(0, 255, 0.25)`, 128, 0, 0, 255)
	assertColor(t, `hsl(0, 255, 0.125)`, 64, 0, 0, 255)
	assertColor(t, `hsl(0, 255, 0.0625)`, 33, 0, 0, 255)
	assertColor(t, `hsl(0, 255, 0.03125)`, 15, 0, 0, 255)
	assertColor(t, `hsl(0, 255, 0.015625)`, 8, 0, 0, 255)

	assertColor(t, `hsl(0, 1.0, 0.5)`, 255, 0, 0, 255)
	assertColor(t, `hsl(0, 1.0, 0.25)`, 128, 0, 0, 255)
	assertColor(t, `hsl(0, 1.0, 0.125)`, 64, 0, 0, 255)
	assertColor(t, `hsl(0, 1.0, 0.0625)`, 33, 0, 0, 255)
	assertColor(t, `hsl(0, 1.0, 0.03125)`, 15, 0, 0, 255)
	assertColor(t, `hsl(0, 1.0, 0.015625)`, 8, 0, 0, 255)

	assertColor(t, `hsv(0, 100%, 100%)`, 255, 0, 0, 255)
	assertColor(t, `hsv(0, 100%, 50%)`, 128, 0, 0, 255)
	assertColor(t, `hsv(0, 100%, 25%)`, 64, 0, 0, 255)
	assertColor(t, `hsv(0, 100%, 12.5%)`, 33, 0, 0, 255)
	assertColor(t, `hsv(0, 100%, 6.25%)`, 15, 0, 0, 255)
	assertColor(t, `hsv(0, 100%, 3.125%)`, 8, 0, 0, 255)
	assertColor(t, `hsv(0, 100%, 1.5625%)`, 5, 0, 0, 255)
}

func TestEquals(t *testing.T) {
	var assert = require.New(t)

	assert.True(Equals(`#2bF1c9`, `#2BF1C9`))
	assert.True(Equals(`#FFFFFF`, `#FFFFFF`))
	assert.True(Equals(`#FFFFFF`, `rgb(255,255,255)`))
	assert.True(Equals(`#FFFFFF80`, `rgba(255,255,255,0.5)`))
	assert.True(Equals(`#00AA00`, `rgb(0,170,0)`))
}

func TestLightenDarken(t *testing.T) {
	var assert = assert.New(t)

	var color Color
	var err error

	// Darken(#800000, 100%) => #000000
	color, err = Darken(`#800000`, 100)
	assert.NoError(err)
	assert.True(color.Equals(`#000000`), fmt.Sprintf("%v", color))

	// Darken(#800000, 20%) => #200000
	color, err = Darken(`#800000`, 20)
	assert.NoError(err)
	assert.True(color.Equals(`#1A0000`), fmt.Sprintf("%v", color))

	// Lighten(#800000, 20%) => #E00000
	color, err = Lighten(`#800000`, 20)
	assert.NoError(err)
	assert.True(color.Equals(`#E60000`), fmt.Sprintf("%v", color))
}

func TestSaturateDesaturate(t *testing.T) {
	var assert = assert.New(t)

	var color Color
	var err error

	// Desaturate(#800000, 100%) => #000000
	color, err = Desaturate(`#800000`, 100)
	assert.NoError(err)
	assert.True(color.Equals(`#404040`), fmt.Sprintf("%v", color))

	// Desaturate(#800000, 20%) => #200000
	color, err = Desaturate(`#ad4038`, 20)
	assert.NoError(err)
	assert.True(color.Equals(`#96544f`), fmt.Sprintf("%v", color))

	// Saturate(#800000, 20%) => #E00000
	color, err = Saturate(`#ad4038`, 20)
	assert.NoError(err)
	assert.True(color.Equals(`#c42b21`), fmt.Sprintf("%v", color))
}

func TestMix(t *testing.T) {
	var assert = assert.New(t)

	var mixed Color
	var err error

	mixed, err = MixN(`#ad4038`, `#0000ff`, 0.8)
	assert.NoError(err)
	assert.True(mixed.Equals(`#8a3360`), fmt.Sprintf("%v", mixed))
}

func TestColorStringers(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(`#FF00CC`, MustParse(`#FF00CC`).String())
	assert.Equal(`#FF00CC`, MustParse(`rgb(255,0,204)`).String())
	assert.Equal(`#FF00CC`, MustParse(`hsl(312, 100%, 50%)`).String())

	assert.Equal(`rgb(255, 0, 204)`, MustParse(`#FF00CC`).StringRGBA())
	assert.Equal(`rgb(255, 0, 204)`, MustParse(`rgb(255,0,204)`).StringRGBA())
	assert.Equal(`rgb(255, 0, 204)`, MustParse(`hsl(312, 100%, 50%)`).StringRGBA())

	assert.Equal(`rgb(255, 0, 204)`, MustParse(`#FF00CCFF`).StringRGBA())
	assert.Equal(`rgba(255, 0, 204, 128)`, MustParse(`#FF00CC80`).StringRGBA())
	assert.Equal(`rgba(255, 0, 204, 64)`, MustParse(`#FF00CC40`).StringRGBA())
	assert.Equal(`rgba(255, 0, 204, 1)`, MustParse(`#FF00CC01`).StringRGBA())

	assert.True(`hsl(312, 100%, 50%)` == MustParse(`#FF00CCFF`).StringHSLA())
	assert.True(`hsla(312, 100%, 50%, 0.5)` == MustParse(`#FF00CC80`).StringHSLA())
}

func TestAdjustHue(t *testing.T) {
	adjustHue(t, `#FF0000`, 120, `#00FF00`)
	adjustHue(t, `#FF0000`, 240, `#0000FF`)
	adjustHue(t, `#ad4038`, 20, `#ad6638`)
}

func TestIsZero(t *testing.T) {
	var c Color

	require.True(t, c.IsZero())

	c = MustParse(`ff00cc`)

	require.False(t, c.IsZero())
}

func TestColorNames(t *testing.T) {
	require.True(t, MustParse(`red`).Equals(`#ff0000`))
	require.True(t, MustParse(`green`).Equals(`#008000`))
	require.True(t, MustParse(`blue`).Equals(`#0000ff`))

	ColorNames[`_testing`] = `#123456`
	require.True(t, MustParse(`_testing`).Equals(`#123456`))
}
