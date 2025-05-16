package geoutil

import (
	"encoding/json"
	"fmt"

	"github.com/ghetzel/go-stockutil/typeutil"
)

var SpeedDisplayUnit = MeasurementSystem(Imperial)

type Speed float64

const (
	MetersPerSecond = 1
	KPH             = 0.277778
	FeetPerSecond   = 0.3048
	MPH             = 0.44704
	Mach            = 340.29
)

func (self Speed) SlowerThan(other Speed) bool {
	return (self < other)
}

func (self Speed) FasterThan(other Speed) bool {
	return (self > other)
}

func (self Speed) Equal(other Speed) bool {
	return (self == other)
}

func (self Speed) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`value`:   float64(self),
		`display`: self.String(),
	})
}

func (self *Speed) UnmarshalJSON(data []byte) error {
	if typeutil.IsScalar(string(data)) {
		var v float64

		if err := json.Unmarshal(data, &v); err == nil {
			*self = Speed(v)
			return nil
		} else {
			return err
		}
	} else {
		var in map[string]any

		if err := json.Unmarshal(data, &in); err == nil {
			*self = Speed(typeutil.V(in[`value`]).Float())
			return nil
		} else {
			return err
		}
	}
}

func (self Speed) String() string {
	switch SpeedDisplayUnit {
	case Metric:
		switch {
		case self >= KPH:
			return fmt.Sprintf("%.0f kph", self/KPH)

		case self >= (Mach * 0.97):
			return fmt.Sprintf("Mach %1.2f", self/Mach)
		}

	case Imperial:
		switch {
		case self < MPH:
			return fmt.Sprintf("%.0f f/s", self/FeetPerSecond)

		case self >= (Mach * 0.97):
			return fmt.Sprintf("Mach %1.2f", self/Mach)

		default:
			return fmt.Sprintf("%.0f mph", self/MPH)
		}
	}

	return fmt.Sprintf("%.0f m/s", self)
}
