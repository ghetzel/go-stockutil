package colorutil

import (
	"fmt"
	"testing"

	"github.com/ghetzel/testify/require"
)

// This benchmark is here both to validate the speed, but also to ensure none of this
// wild floating point math accumulates any noticable errors.
func BenchmarkKelvinToColor(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for k, rgba := range map[int][3]uint8{
			1850: [3]uint8{255, 130, 0},
			1999: [3]uint8{255, 139, 0},
			2000: [3]uint8{255, 139, 0},
			2500: [3]uint8{255, 162, 71},
			3000: [3]uint8{255, 180, 108},
			3500: [3]uint8{255, 195, 138},
			4000: [3]uint8{255, 208, 164},
			4500: [3]uint8{255, 219, 185},
			5000: [3]uint8{255, 228, 205},
			5500: [3]uint8{255, 236, 223},
			6000: [3]uint8{255, 243, 239},
			6500: [3]uint8{255, 250, 254},
			6599: [3]uint8{255, 251, 255},
			6600: [3]uint8{255, 249, 255},
			7000: [3]uint8{245, 243, 255},
			7500: [3]uint8{234, 237, 255},
			8000: [3]uint8{225, 232, 255},
			8500: [3]uint8{218, 228, 255},
			9000: [3]uint8{213, 225, 255},
			9001: [3]uint8{213, 225, 255},
		} {
			var r, g, b uint8

			r, g, b, _ = KelvinToColor(k).RGBA255()

			if r != rgba[0] {
				panic(fmt.Sprintf("bad value for r:%dK; expected %d, got %d", k, r, rgba[0]))
			}

			if g != rgba[1] {
				panic(fmt.Sprintf("bad value for g:%dK; expected %d, got %d", k, g, rgba[1]))
			}

			if b != rgba[2] {
				panic(fmt.Sprintf("bad value for b:%dK; expected %d, got %d", k, b, rgba[2]))
			}
		}
	}
}

func TestKelvinToColor(t *testing.T) {
	var assert = require.New(t)

	for k, rgba := range map[int][3]uint8{
		1850: [3]uint8{255, 130, 0},
		1999: [3]uint8{255, 139, 0},
		2000: [3]uint8{255, 139, 0},
		2500: [3]uint8{255, 162, 71},
		3000: [3]uint8{255, 180, 108},
		3500: [3]uint8{255, 195, 138},
		4000: [3]uint8{255, 208, 164},
		4500: [3]uint8{255, 219, 185},
		5000: [3]uint8{255, 228, 205},
		5500: [3]uint8{255, 236, 223},
		6000: [3]uint8{255, 243, 239},
		6500: [3]uint8{255, 250, 254},
		6599: [3]uint8{255, 251, 255},
		6600: [3]uint8{255, 249, 255},
		7000: [3]uint8{245, 243, 255},
		7500: [3]uint8{234, 237, 255},
		8000: [3]uint8{225, 232, 255},
		8500: [3]uint8{218, 228, 255},
		9000: [3]uint8{213, 225, 255},
		9001: [3]uint8{213, 225, 255},
	} {
		var r, g, b, a uint8

		r, g, b, a = KelvinToColor(k).RGBA255()
		assert.EqualValues(r, rgba[0], fmt.Sprintf("r:%dK", k))
		assert.EqualValues(g, rgba[1], fmt.Sprintf("g:%dK", k))
		assert.EqualValues(b, rgba[2], fmt.Sprintf("b:%dK", k))
		assert.EqualValues(a, 255, fmt.Sprintf("a:%dK", k))
	}
}
