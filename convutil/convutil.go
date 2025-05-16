package convutil

import (
	"fmt"

	"github.com/ghetzel/go-stockutil/mathutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/martinlindhe/unit"
)

var ConvertRoundToPlaces = 6

func parseUnit(in any) Unit {
	var out = Invalid

	if v, ok := in.(Unit); ok {
		out = v
	} else if vS, ok := in.(string); ok {
		if v, err := ParseUnit(vS); err == nil {
			out = v
		} else {
			panic(err.Error())
		}
	} else if vS, ok := in.(fmt.Stringer); ok {
		if v, err := ParseUnit(vS.String()); err == nil {
			out = v
		} else {
			panic(err.Error())
		}
	}

	return out
}

func MustConvert(in any, from any, to any) float64 {
	if value, err := Convert(in, from, to); err == nil {
		return value
	} else {
		panic(err.Error())
	}
}

func Convert(in any, from any, to any) (float64, error) {
	if v, err := ExactConvert(in, from, to); err == nil {
		return mathutil.RoundPlaces(v, ConvertRoundToPlaces), nil
	} else {
		return 0, err
	}
}

func MustExactConvert(in any, from any, to any) float64 {
	if value, err := Convert(in, from, to); err == nil {
		return value
	} else {
		panic(err.Error())
	}
}

func ExactConvert(in any, from any, to any) (float64, error) {
	var fromU = parseUnit(from)
	var toU = parseUnit(to)

	if !fromU.IsValid() {
		return 0, fmt.Errorf("invalid 'from' unit")
	}

	if !toU.IsValid() {
		return 0, fmt.Errorf("invalid 'to' unit")
	}

	if v, err := stringutil.ConvertToFloat(in); err == nil {
		if from == to {
			return v, nil
		}

		if fromU.Family() != toU.Family() {
			return 0, fmt.Errorf("units '%v' and '%v' are not convertible to each other", fromU, toU)
		}

		var baseConvert = func() any {
			switch fromU.Family() {
			case TemperatureUnits:
				return convertTemperature(v, fromU, toU)
			case SpeedUnits:
				return convertSpeed(v, fromU, toU)
			case LengthUnits:
				return convertDistance(v, fromU, toU)
			}

			return fmt.Errorf("cannot convert from %v to %v", from, to)
		}

		var converted = baseConvert()

		if vF, ok := converted.(float64); ok {
			return vF, nil
		} else if err, ok := converted.(error); ok {
			return 0, err
		} else {
			return 0, fmt.Errorf("unspecified error")
		}
	} else {
		return 0, err
	}
}

func convertTemperature(v float64, from Unit, to Unit) float64 {
	var c unit.Temperature

	switch from {
	case Celsius:
		c = unit.FromCelsius(v)
	case Fahrenheit:
		c = unit.FromFahrenheit(v)
	case Kelvin:
		c = unit.FromKelvin(v)
	default:
		return 0
	}

	switch to {
	case Fahrenheit:
		return c.Fahrenheit()
	case Celsius:
		return c.Celsius()
	case Kelvin:
		return c.Kelvin()
	}

	return 0
}

func convertSpeed(v float64, from Unit, to Unit) float64 {
	var c = unit.Speed(v)

	switch from {
	case MilesPerHour:
		c *= unit.MilesPerHour
	case KilometersPerHour:
		c *= unit.KilometersPerHour
	}

	switch to {
	case MetersPerSecond:
		return c.MetersPerSecond()
	case MilesPerHour:
		return c.MilesPerHour()
	case KilometersPerHour:
		return c.KilometersPerHour()
	}

	return 0
}

func convertDistance(v float64, from Unit, to Unit) float64 {
	var c = unit.Length(v)

	switch from {
	case Feet:
		c *= unit.Foot
	case Miles:
		c *= unit.Mile
	case NauticalMiles:
		c *= unit.NauticalMile
	case AU:
		c *= unit.AstronomicalUnit
	case Lightyears, Lightminutes, Lightseconds:
		c *= unit.LightYear
	}

	switch to {
	case Meters:
		return c.Meters()
	case Feet:
		return c.Feet()
	case Miles:
		return c.Miles()
	case NauticalMiles:
		return c.NauticalMiles()
	case AU:
		return c.AstronomicalUnits()
	case Lightyears:
		return c.LightYears()
	case Lightminutes:
		return (c.LightYears() * 525600)
	case Lightseconds:
		return (c.LightYears() * 31536000)
	}

	return 0
}
