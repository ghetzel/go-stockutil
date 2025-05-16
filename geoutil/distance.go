package geoutil

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var rxDistanceExtract = regexp.MustCompile(`^(?P<number>\d+(?:\.\d+)?)\s*(?P<unit>[\w\s]+)\W*$`)

type Distance float64

const (
	Meter        = 1
	Kilometer    = 1000
	Foot         = 0.3048
	Yard         = 0.9144
	Mile         = 1609.344
	NauticalMile = 1852
)

var DistanceDisplayUnit = MeasurementSystem(Imperial)

func MustParseDistance(in any) Distance {
	if distance, err := ParseDistance(in); err == nil {
		return distance
	} else {
		panic(`invalid distance: ` + err.Error())
	}
}

func ParseDistance(in any) (Distance, error) {
	if typeutil.IsZero(in) {
		return 0, nil
	}

	if match := rxutil.Match(rxDistanceExtract, strings.TrimSpace(fmt.Sprintf("%v", in))); match != nil {
		if v := typeutil.V(match.Group(`number`)).Float(); v >= 0 {
			var unit = match.Group(`unit`)
			unit = strings.TrimSpace(unit)
			unit = strings.ToLower(unit)
			unit = strings.TrimSuffix(unit, `s`)

			switch unit {
			case `meter`, `m`:
				return Distance(v), nil
			case `kilometer`, `km`:
				return Distance(v) * Kilometer, nil
			case `mile`, `mi`:
				return Distance(v) * Mile, nil
			case `feet`, `foot`, `ft`:
				return Distance(v) * Foot, nil
			case `yard`, `yd`:
				return Distance(v) * Yard, nil
			case `nm`, `nautical mile`:
				return Distance(v) * NauticalMile, nil
			default:
				return 0, fmt.Errorf("Unrecognized distance unit %q", unit)
			}
		} else {
			return 0, fmt.Errorf("Unable to extract number from distance value")
		}
	} else if v := typeutil.V(in).Float(); v >= 0 {
		return Distance(v), nil
	} else {
		return 0, fmt.Errorf("unable to parse distance value")
	}
}

func (self Distance) Within(other Distance) bool {
	return (self <= other)
}

func (self Distance) Beyond(other Distance) bool {
	return (self > other)
}

func (self Distance) Equal(other Distance) bool {
	return (self == other)
}

func (self Distance) MetricString() string {
	switch {
	case self >= 10*Kilometer:
		return fmt.Sprintf("%.0f kilometers", self/Kilometer)

	case self > Kilometer:
		return fmt.Sprintf("%.1f kilometers", self/Kilometer)

	case self == Kilometer:
		return fmt.Sprintf("%.0f kilometer", self/Kilometer)

	default:
		return fmt.Sprintf("%.0f meters", self)

	}
}

func (self Distance) ImperialString() string {
	switch {
	case self >= 5*Mile:
		return fmt.Sprintf("%.0f miles", self/Mile)

	case self >= (0.9*Mile) && self <= (1.1*Mile):
		return fmt.Sprintf("%.0f mile", self/Mile)

	case self >= (0.1 * Mile):
		return fmt.Sprintf("%.1f miles", self/Mile)

	default:
		return fmt.Sprintf("%.0f feet", self/Foot)
	}
}

func (self Distance) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`value`: float64(self),
		`units`: map[string]any{
			`default`:  self.String(),
			`imperial`: self.ImperialString(),
			`metric`:   self.MetricString(),
		},
	})
}

func (self *Distance) UnmarshalJSON(data []byte) error {
	if typeutil.IsScalar(string(data)) {
		var v float64

		if err := json.Unmarshal(data, &v); err == nil {
			*self = Distance(v)
			return nil
		} else {
			return err
		}
	} else {
		var in map[string]any

		if err := json.Unmarshal(data, &in); err == nil {
			*self = Distance(typeutil.V(in[`value`]).Float())
			return nil
		} else {
			return err
		}
	}
}

func (self Distance) String() string {
	switch DistanceDisplayUnit {
	case Metric:
		return self.MetricString()

	case Imperial:
		return self.ImperialString()
	}

	return fmt.Sprintf("%f meters", self)
}
