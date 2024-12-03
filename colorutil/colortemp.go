package colorutil

import "math"

// These values were derived by Neil Bartlett (c) 2019
//
//	https://www.zombieprototypes.com/?p=210
//	see also: https://github.com/neilbartlett/color-temperature
//
// Based on original work by Tanner Helland (c) 2014
//
//	http://www.tannerhelland.com/4435/convert-temperature-rgb-algorithm-code/
var hellandBartlettOverEq6600K_RedA float64 = 351.97690566805693
var hellandBartlettOverEq6600K_RedB float64 = 0.114206453784165
var hellandBartlettOverEq6600K_RedC float64 = -40.25366309332127
var hellandBartlettOverEq6600K_RedX float64 = -55
var hellandBartlettUnder6600K_GreenA float64 = -155.25485562709179
var hellandBartlettUnder6600K_GreenB float64 = -0.44596950469579133
var hellandBartlettUnder6600K_GreenC float64 = 104.49216199393888
var hellandBartlettUnder6600K_GreenX float64 = -2
var hellandBartlettOverEq6600K_GreenA float64 = 325.4494125711974
var hellandBartlettOverEq6600K_GreenB float64 = 0.07943456536662342
var hellandBartlettOverEq6600K_GreenC float64 = -28.0852963507957
var hellandBartlettOverEq6600K_GreenX float64 = -50
var hellandBartlettOver2000K_Under6600K_BlueA float64 = -254.76935184120902
var hellandBartlettOver2000K_Under6600K_BlueB float64 = 0.8274096064007395
var hellandBartlettOver2000K_Under6600K_BlueC float64 = 115.67994401066147
var hellandBartlettOver2000K_Under6600K_BlueX float64 = -10

// Takes a color temperature in degrees Kelvin and returns a valid Color for that temperature.
// This function works best between 1000K and 40000K.
func KelvinToColor(kelvin int) (color Color) {
	color.a = 1.0

	// TODO: Start with a temperature, in Kelvin, somewhere between 1000 and 40000.

	var rA, rB, rC, rX float64
	var g1A, g1B, g1C, g1X float64
	var g2A, g2B, g2C, g2X float64
	var bA, bB, bC, bX float64

	rA = hellandBartlettOverEq6600K_RedA
	rB = hellandBartlettOverEq6600K_RedB
	rC = hellandBartlettOverEq6600K_RedC
	rX = hellandBartlettOverEq6600K_RedX
	g1A = hellandBartlettUnder6600K_GreenA
	g1B = hellandBartlettUnder6600K_GreenB
	g1C = hellandBartlettUnder6600K_GreenC
	g1X = hellandBartlettUnder6600K_GreenX
	g2A = hellandBartlettOverEq6600K_GreenA
	g2B = hellandBartlettOverEq6600K_GreenB
	g2C = hellandBartlettOverEq6600K_GreenC
	g2X = hellandBartlettOverEq6600K_GreenX
	bA = hellandBartlettOver2000K_Under6600K_BlueA
	bB = hellandBartlettOver2000K_Under6600K_BlueB
	bC = hellandBartlettOver2000K_Under6600K_BlueC
	bX = hellandBartlettOver2000K_Under6600K_BlueX

	// scale the value to make working with some of this easier
	temp := float64(kelvin) / 100.0

	// Calculate Red
	if temp < 66 {
		color.r = 1.0
	} else {
		// a + b*x + c*ln(x)
		x := (temp + rX)
		r255 := (rA + (rB * x) + (rC * math.Log(x)))

		if r255 < 0 {
			r255 = 0
		} else if r255 > 255 {
			r255 = 255
		}

		color.r = r255 / 255
	}

	// Calculate Green
	if temp < 66 {
		x := (temp + g1X)
		g255 := (g1A + (g1B * x) + (g1C * math.Log(x)))

		if g255 < 0 {
			g255 = 0
		} else if g255 > 255 {
			g255 = 255
		}

		color.g = g255 / 255
	} else {
		x := (temp + g2X)
		g255 := (g2A + (g2B * x) + (g2C * math.Log(x)))

		if g255 < 0 {
			g255 = 0
		} else if g255 > 255 {
			g255 = 255
		}

		color.g = g255 / 255
	}

	// Calculate Blue
	if temp <= 20 {
		color.b = 0
	} else if temp >= 66 {
		color.b = 1.0
	} else {
		x := (temp + bX)
		b255 := (bA + (bB * x) + (bC * math.Log(x)))

		if b255 < 0 {
			b255 = 0
		} else if b255 > 255 {
			b255 = 255
		}

		color.b = b255 / 255
	}

	return
}
