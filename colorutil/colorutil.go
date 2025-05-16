// Utilities for parsing and manipulating colors.
package colorutil

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/ghetzel/go-stockutil/mathutil"
	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"k8s.io/apimachinery/pkg/util/json"
)

type adjustment int

const (
	hue adjustment = iota
	saturation
	lightness
)

var rgbaPattern = `rgba?\((?P<r>\d+(?:\.\d+)?%?)\s*,\s*(?P<g>\d+(?:\.\d+)?%?)\s*,\s*(?P<b>\d+(?:\.\d+)?%?)\s*(?:,\s*(?P<a>\d+(?:\.\d+)?%?)?\s*)?\)`
var hslaPattern = `hs(?P<LorV>[lv])a?\((?P<h>\d+(?:\.\d+)?%?)(?:deg)?\s*,\s*(?P<s>\d+(?:\.\d+)?%?)\s*,\s*(?P<lv>\d+(?:\.\d+)?%?)\s*(?:,\s*(?P<a>\d+(?:\.\d+)?%?)?\s*)?\)`
var hexPattern = `#?(?P<r>[0-9a-fA-F]{2})(?P<g>[0-9a-fA-F]{2})(?P<b>[0-9a-fA-F]{2})(?P<a>[0-9a-fA-F]{2})?`

const precision = 2

type lmodel int

const (
	Invalid lmodel = iota
	HSI
	HSL
	HSV
)

type Color struct {
	r float64
	g float64
	b float64
	a float64
}

func (self Color) IsZero() bool {
	return self.r == 0 && self.g == 0 && self.b == 0 && self.a == 0
}

func (self Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.String())
}

func (self Color) String() string {
	r, g, b, a := self.RGBA255()

	if self.a < 1 {
		return fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a)
	} else {
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
}

func (self Color) StringRGBA() string {
	r, g, b, a := self.RGBA255()

	if self.a < 1 {
		return fmt.Sprintf(
			"rgba(%d, %d, %d, %d)",
			r,
			g,
			b,
			a,
		)
	} else {
		return fmt.Sprintf(
			"rgb(%d, %d, %d)",
			r,
			g,
			b,
		)
	}
}

func (self Color) StringHSLA() string {
	h, s, l := self.HSL()

	if self.a < 1 {
		return fmt.Sprintf(
			`hsla(%d, %d%%, %d%%, %g)`,
			int(math.Mod(mathutil.Round(h), 360)),
			int(s*100.0),
			int(l*100.0),
			mathutil.RoundPlaces(self.a, 2),
		)
	} else {
		return fmt.Sprintf(
			"hsl(%d, %d%%, %d%%)",
			int(math.Mod(mathutil.Round(h), 360)),
			int(s*100.0),
			int(l*100.0),
		)
	}
}

// Return the current color as 4x 8-bit RGB values, each [0, 255].
func (self Color) RGBA255() (uint8, uint8, uint8, uint8) {
	var r, g, b, a uint8

	r = uint8(self.r*255.0 + 0.5)
	g = uint8(self.g*255.0 + 0.5)
	b = uint8(self.b*255.0 + 0.5)
	a = uint8(self.a*255.0 + 0.5)

	return r, g, b, a
}

// Return the current color as a 32-bit uint quad, implementing the color.Color interface.
func (self Color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(self.r * float64(0xFFFF)),
		uint32(self.g * float64(0xFFFF)),
		uint32(self.b * float64(0xFFFF)),
		uint32(self.a * float64(0xFFFF))
}

// Return the current color as hue (out of 360°), saturation [0, 1], and lightness [0, 1].
func (self Color) HSL() (float64, float64, float64) {
	return rgb2lhs(HSL, self.r, self.g, self.b)
}

// Return the current color as hue (out of 360°), saturation [0, 1], and value [0, 1].
func (self Color) HSV() (float64, float64, float64) {
	return rgb2lhs(HSV, self.r, self.g, self.b)
}

// Return the current color as hue (out of 360°), saturation [0, 1], and intensity [0, 1].
func (self Color) HSI() (float64, float64, float64) {
	return rgb2lhs(HSI, self.r, self.g, self.b)
}

// Return the color as a color.RGBA struct.
func (self Color) NativeRGBA() (c color.RGBA) {
	c.R, c.G, c.B, c.A = self.RGBA255()
	return
}

// Return whether the given color is equal to this one in the 24-bit RGB (RGB255) color space
func (self Color) Equals(other any) bool {
	if otherC, err := Parse(other); err == nil {
		r1, g1, b1, a1 := self.RGBA255()
		r2, g2, b2, a2 := otherC.RGBA255()

		if r1 == r2 {
			if g1 == g2 {
				if b1 == b2 {
					if a1 == a2 {
						return true
					}
				}
			}
		}
	}

	return false
}

// Parse the given value into a Color or panic.
func MustParse(value any) Color {
	if color, err := Parse(value); err == nil {
		return color
	} else {
		panic(err)
	}
}

// Parse the given value into a Color or return an error.
func Parse(value any) (Color, error) {
	var colorC Color

	colorC.a = 1

	if c, ok := value.(Color); ok {
		return c, nil
	} else {
		var colorS = fmt.Sprintf("%v", value)
		colorS = strings.TrimSpace(colorS)
		colorS = strings.ToLower(colorS)

		if c, ok := ColorNames[colorS]; ok && strings.TrimSpace(c) != `` {
			colorS = c
		}

		if rgba := rxutil.Match(rgbaPattern, colorS); rgba != nil {
			// handles rgb(), rgba() patterns
			for v, factorS := range rgba.NamedCaptures() {
				var denom float64 = 0xFF

				if strings.HasSuffix(factorS, `%`) {
					denom = 100
					factorS = factorS[:len(factorS)-1]
				}

				if factorF, err := stringutil.ConvertToFloat(factorS); err == nil {
					if v != `a` && factorF > 1 {
						factorF = (factorF / denom)
					}

					switch v {
					case `r`:
						colorC.r = factorF
					case `g`:
						colorC.g = factorF
					case `b`:
						colorC.b = factorF
					case `a`:
						if factorS != `` {
							if factorF > 1 {
								factorF = (factorF / denom)
							}

							colorC.a = factorF
						}
					}
				} else {
					return Color{}, fmt.Errorf("Invalid rgba() value in '%v' position: %v", v, err)
				}
			}

		} else if hsla := rxutil.Match(hslaPattern, colorS); hsla != nil {
			// handles hsl(), hsla(), hsv(), hsva() patterns
			var h, s, lv float64
			var isHSV bool

			if hsla.Group(`LorV`) == `v` {
				isHSV = true
			}

			for v, factorS := range hsla.NamedCaptures() {
				if v == `LorV` {
					continue
				}

				var denom float64 = 1

				if strings.HasSuffix(factorS, `%`) {
					denom = 100
					factorS = strings.TrimSuffix(factorS, `%`)
				}

				if factorF, err := stringutil.ConvertToFloat(factorS); err == nil {
					if factorF > 1 && denom == 1 {
						denom = 0xFF
					}

					switch v {
					case `h`:
						h = factorF
					case `s`:
						s = (factorF / denom)
					case `lv`:
						lv = (factorF / denom)
					case `a`:
						if factorS != `` {
							if factorF > 1 {
								factorF = (factorF / denom)
							}

							colorC.a = factorF
						}
					}
				} else {
					return Color{}, fmt.Errorf("Invalid rgba() value in '%v' position: %v", v, err)
				}
			}

			if isHSV {
				colorC.r, colorC.g, colorC.b = hsv2rgb(h, s, lv)
			} else {
				colorC.r, colorC.g, colorC.b = hsl2rgb(h, s, lv)
			}
		} else if hexa := rxutil.Match(hexPattern, colorS); hexa != nil {
			for v, hexbyte := range hexa.NamedCaptures() {
				if hexbyte != `` {
					if vB, err := hex.DecodeString(hexbyte); err == nil {
						if len(vB) == 1 {
							switch v {
							case `r`:
								colorC.r = (float64(vB[0]) / 0xFF)
							case `g`:
								colorC.g = (float64(vB[0]) / 0xFF)
							case `b`:
								colorC.b = (float64(vB[0]) / 0xFF)
							case `a`:
								colorC.a = (float64(vB[0]) / 0xFF)
							}
						} else {
							return Color{}, fmt.Errorf("Invalid hex byte '%v'", vB)
						}
					} else {
						return Color{}, fmt.Errorf("Invalid hex value '%v'", hexbyte)
					}
				}
			}
		}
	}

	return colorC, nil
}

// Return whether two colors are equivalent in the 24-bit RGB (RGB255) color space.
func Equals(first any, second any) bool {
	if first != nil && second != nil {
		if firstC, err := Parse(first); err == nil {
			return firstC.Equals(second)
		}
	}

	return false
}

// adjust the a property of given color by a specified factor:
//
//		Hue: add/subtract the given number of degrees on the HSL color wheel
//	 Saturation: saturate/desaturate the color by the given amount
//	 Lightness: lighten/darken the color by the given amount
func adjust(what adjustment, in any, factor float64) (Color, error) {
	if sample, err := Parse(in); err == nil {
		h, s, l := sample.HSL()

		switch what {
		case lightness:
			l += factor

			if l < 0 {
				l = 0
			} else if l > 1 {
				l = 1
			}
		case saturation:
			s += factor

			if s < 0 {
				s = 0
			} else if s > 1 {
				s = 1
			}
		case hue:
			h += factor

		default:
			return Color{}, fmt.Errorf("invalid what value")
		}

		sample.r, sample.g, sample.b = hsl2rgb(h, s, l)

		return sample, nil
	} else {
		return Color{}, err
	}
}

// Darken the given color by a certain percent.  Consistent with the results of the
// Sass darken() function.
func Darken(in any, percent int) (Color, error) {
	return adjust(lightness, in, -1*(float64(percent)/100.0))
}

// Lighten the given color by a certain percent.  Consistent with the results of the
// Sass lighten() function.
func Lighten(in any, percent int) (Color, error) {
	return adjust(lightness, in, float64(percent)/100.0)
}

// Saturate the given color by a certain percent.  Consistent with the results of the
// Sass saturate() function.
func Saturate(in any, percent int) (Color, error) {
	return adjust(saturation, in, float64(percent)/100.0)
}

// Desaturate the given color by a certain percent.  Consistent with the results of the
// Sass desaturate() function.
func Desaturate(in any, percent int) (Color, error) {
	return adjust(saturation, in, -1*(float64(percent)/100.0))
}

// Adjust the hue of the given color by the specified number of degrees.
func AdjustHue(in any, degrees float64) (Color, error) {
	return adjust(hue, in, degrees)
}

// Mix two colors, producing a third.  The weight value specifies how much of the first color
// should be included.  Consistent with the results of the Sass mix() function.
func MixN(first any, second any, weight float64) (Color, error) {
	if firstC, err := Parse(first); err == nil {
		if secondC, err := Parse(second); err == nil {
			var factor = 2*weight - 1
			var alpha = firstC.a - secondC.a

			var weight1, weight2 float64

			if factor*alpha == -1 {
				weight1 = factor
			} else {
				weight1 = ((factor + alpha) / (1 + factor*alpha)) + 1
			}

			weight1 = weight1 / 2
			weight2 = 1 - weight1

			var mixed = Color{
				r: (weight1 * firstC.r) + (weight2 * secondC.r),
				g: (weight1 * firstC.g) + (weight2 * secondC.g),
				b: (weight1 * firstC.b) + (weight2 * secondC.b),
				a: (firstC.a * weight) + (secondC.a * (1 - weight)),
			}

			return mixed, nil
		} else {
			return Color{}, fmt.Errorf("second color: %v", err)
		}
	} else {
		return Color{}, fmt.Errorf("first color: %v", err)
	}
}

// Mix two colors in equal parts, producing a third.
func Mix(first any, second any) (Color, error) {
	return MixN(first, second, 0.5)
}

// Given HSL values (where hue is given in degrees (out of 360°), saturation
// and lightness are [0, 1]), return the corresponding RGB values (where each value
// is [0, 255]).
func HslToRgb(hue float64, saturation float64, lightness float64) (uint8, uint8, uint8) {
	r, g, b := hsl2rgb(hue, saturation, lightness)

	return uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0)
}

// Given HSL values (where hue is given in degrees (out of 360°), saturation
// and value are [0, 1]), return the corresponding RGB values (where each value
// is [0, 255]).
func HsvToRgb(hue float64, saturation float64, value float64) (uint8, uint8, uint8) {
	r, g, b := hsv2rgb(hue, saturation, value)

	return uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0)
}

// Given RGB values (where each value is [0, 255]), return the hue (in degrees), saturation,
// and lightness (where each is [0, 1]).
func RgbToHsl(r float64, g float64, b float64) (float64, float64, float64) {
	return rgb2lhs(HSL, r/255, g/255, b/255)
}

// Internal HSL->RGB function for doing conversions using float inputs (saturation, lightness) and
// outputs (for R, G, and B).
func hsl2rgb(hueDegrees float64, saturation float64, lightness float64) (float64, float64, float64) {
	return hs2rgb(false, hueDegrees, saturation, lightness)
}

// Internal HSV->RGB function for doing conversions using float inputs (saturation, value) and
// outputs (for R, G, and B).
func hsv2rgb(hueDegrees float64, saturation float64, value float64) (float64, float64, float64) {
	return hs2rgb(true, hueDegrees, saturation, value)
}

// Internal generic implementation for converting HSV, HSL to RGB.
// Basically a direct implementation of this: https://en.wikipedia.org/wiki/HSL_and_HSV#Converting_to_RGB
func hs2rgb(isValue bool, hueDegrees float64, saturation float64, lightOrVal float64) (float64, float64, float64) {
	var r, g, b float64

	hueDegrees = math.Mod(hueDegrees, 360)

	if saturation == 0 {
		r = lightOrVal
		g = lightOrVal
		b = lightOrVal
	} else {
		var chroma float64
		var m float64

		if isValue {
			chroma = lightOrVal * saturation
		} else {
			chroma = (1 - math.Abs((2*lightOrVal)-1)) * saturation
		}

		var hueSector = hueDegrees / 60

		var intermediate = chroma * (1 - math.Abs(
			math.Mod(hueSector, 2)-1,
		))

		switch {
		case hueSector >= 0 && hueSector <= 1:
			r = chroma
			g = intermediate
			b = 0

		case hueSector > 1 && hueSector <= 2:
			r = intermediate
			g = chroma
			b = 0

		case hueSector > 2 && hueSector <= 3:
			r = 0
			g = chroma
			b = intermediate

		case hueSector > 3 && hueSector <= 4:
			r = 0
			g = intermediate
			b = chroma
		case hueSector > 4 && hueSector <= 5:
			r = intermediate
			g = 0
			b = chroma

		case hueSector > 5 && hueSector <= 6:
			r = chroma
			g = 0
			b = intermediate

		default:
			panic(fmt.Errorf("hue input %v yielded sector %v", hueDegrees, hueSector))
		}

		if isValue {
			m = lightOrVal - chroma
		} else {
			m = lightOrVal - (chroma / 2)
		}

		r += m
		g += m
		b += m

	}

	r = mathutil.RoundPlaces(r, precision)
	g = mathutil.RoundPlaces(g, precision)
	b = mathutil.RoundPlaces(b, precision)

	return r, g, b
}

// Internal implementation converting RGB to HSL, HSV, or HSI.
// Basically a direct implementation of this: https://en.wikipedia.org/wiki/HSL_and_HSV#General_approach
func rgb2lhs(lightnessModel lmodel, r float64, g float64, b float64) (float64, float64, float64) {
	var h, s, lvi float64
	var huePrime float64

	// hue
	// ---------------------------------------------------------------------------------------------
	var max = math.Max(math.Max(r, g), b)
	var min = math.Min(math.Min(r, g), b)
	var chroma = (max - min)

	if chroma == 0 {
		h = 0
	} else {
		if r == max {
			huePrime = math.Mod(((g - b) / chroma), 6)
		} else if g == max {
			huePrime = ((b - r) / chroma) + 2

		} else if b == max {
			huePrime = ((r - g) / chroma) + 4

		}

		h = huePrime * 60
	}

	// lightness
	// ---------------------------------------------------------------------------------------------
	if r == g && g == b {
		lvi = r
	} else {
		switch lightnessModel {
		case HSI:
			lvi = (r + b + g) / 3
		case HSV:
			lvi = max
		case HSL:
			lvi = (max + min) / 2
		}
	}

	// saturation
	// ---------------------------------------------------------------------------------------------
	switch lightnessModel {
	case HSI:
		if lvi == 0 {
			s = 0
		} else {
			s = (1 - (min / lvi))
		}

	case HSV:
		if lvi == 0 {
			s = 0
		} else {
			s = (chroma / lvi)
		}

	case HSL:
		if lvi == 1 {
			s = 0
		} else {
			s = (chroma / (1 - math.Abs(2*lvi-1)))
		}
	}

	if math.IsNaN(s) {
		s = 0
	}

	h = mathutil.RoundPlaces(h, precision)
	s = mathutil.RoundPlaces(s, precision)
	lvi = mathutil.RoundPlaces(lvi, precision)

	if h < 0 {
		h = 360 + h
	}

	return h, s, lvi
}
