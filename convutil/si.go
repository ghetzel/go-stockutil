package convutil

import (
	"fmt"

	"go.gary.cool/go-stockutil/mathutil"
	"go.gary.cool/go-stockutil/typeutil"
)

type SIExponents int

const (
	Quetta SIExponents = 30
	Ronna              = 27
	Yotta              = 24
	Zetta              = 21
	Exa                = 18
	Peta               = 15
	Tera               = 12
	Giga               = 9
	Mega               = 6
	Kilo               = 3
	Hecto              = 2
	Deca               = 1
	Deci               = -1
	Centi              = -2
	Milli              = -3
	Micro              = -6
	Nano               = -9
	Pico               = -12
	Femto              = -15
	Atto               = -18
	Zepto              = -21
	Yocto              = -24
	Ronto              = -27
	Quecto             = -30
)

type Bytes float64

const (
	Byte       Bytes = 1
	Kilobyte   Bytes = 1024
	Megabyte   Bytes = 1048576
	Gigabyte   Bytes = 1073741824
	Terabyte   Bytes = 1099511627776
	Petabyte   Bytes = 1125899906842624
	Exabyte    Bytes = 1152921504606846976
	Zettabyte  Bytes = 1180591620717411303424
	Yottabyte  Bytes = 1208925819614629174706176
	Brontobyte Bytes = 1237940039285380274899124224 // unofficial
	Ronnabyte  Bytes = 1237940039285380274899124224
	Quettabyte Bytes = 1267650600228229401496703205376
)

func (self Bytes) To(unit Bytes) string {
	value, suffix := self.Convert(unit)
	return typeutil.String(value) + suffix
}

func (self Bytes) String() string {
	value, suffix := self.Auto()

	return typeutil.String(value) + suffix
}

func (self Bytes) Format(format string, as ...Bytes) string {
	var value float64
	var suffix string

	if len(as) > 0 {
		value, suffix = self.Convert(as[0])
	} else {
		value, suffix = self.Auto()
	}

	return fmt.Sprintf(format, value, suffix)
}

func (self Bytes) Auto() (float64, string) {
	return self.Convert(0)
}

func (self Bytes) Convert(to Bytes) (float64, string) {
	var value float64
	var suffix string

	switch {
	case self >= Quettabyte || to == Quettabyte:
		suffix = `QB`
		value = float64(self / Quettabyte)
	case self >= Ronnabyte || to == Ronnabyte:
		suffix = `RB`
		value = float64(self / Ronnabyte)
	case self >= Yottabyte || to == Yottabyte:
		suffix = `YB`
		value = float64(self / Yottabyte)
	case self >= Zettabyte || to == Zettabyte:
		suffix = `ZB`
		value = float64(self / Zettabyte)
	case self >= Exabyte || to == Exabyte:
		suffix = `EB`
		value = float64(self / Exabyte)
	case self >= Petabyte || to == Petabyte:
		suffix = `PB`
		value = float64(self / Petabyte)
	case self >= Terabyte || to == Terabyte:
		suffix = `TB`
		value = float64(self / Terabyte)
	case self >= Gigabyte || to == Gigabyte:
		suffix = `GB`
		value = float64(self / Gigabyte)
	case self >= Megabyte || to == Megabyte:
		suffix = `MB`
		value = float64(self / Megabyte)
	case self >= Kilobyte || to == Kilobyte:
		suffix = `KB`
		value = float64(self / Kilobyte)
	default:
		suffix = `B`
		value = float64(self)
	}

	value = mathutil.RoundPlaces(value, mathutil.LeadingSignificantZeros(value, 31)+1)

	return value, suffix
}
