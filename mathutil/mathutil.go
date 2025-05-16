package mathutil

import (
	"math"
)

const (
	uvone    = 0x3FF0000000000000
	mask     = 0x7FF
	shift    = 64 - 11 - 1
	bias     = 1023
	signMask = 1 << 63
	fracMask = 1<<shift - 1
)

func Clamp(value float64, lower float64, upper float64) float64 {
	value = ClampLower(value, lower)
	value = ClampUpper(value, upper)
	return value
}

func ClampLower(value float64, lower float64) float64 {
	if value < lower {
		value = lower
	}

	return value
}

func ClampUpper(value float64, upper float64) float64 {
	if value > upper {
		value = upper
	}

	return value
}

func RoundPlaces(x float64, places int) float64 {
	if places <= 0 {
		return Round(x)
	} else {
		var multi = math.Pow(10, Clamp(float64(places), 0, 16))
		return (Round(x*multi) / multi)
	}
}

func LeadingSignificantZeros(x float64, places int) int {
	var zcount int

	for _, div := math.Modf(x); div > 0 && div < 1; div = RoundPlaces(div/0.1, places) {
		zcount += 1
	}

	if zcount > 0 {
		return zcount - 1
	} else {
		return 0
	}
}
