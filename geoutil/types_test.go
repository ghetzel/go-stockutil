package geoutil

import (
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

func TestLocationSpeedFrom(t *testing.T) {
	var assert = require.New(t)

	assert.Zero(NullIsland().SpeedFrom(NullIsland()))

	var from = &Location{
		Latitude:  0,
		Longitude: 1,
		Timestamp: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	var to = &Location{
		Latitude:  0,
		Longitude: 2,
		Timestamp: time.Date(2010, 1, 1, 1, 0, 0, 0, time.UTC),
	}

	// 1deg of latitude over 1 hour ~ 111.19 kph
	assert.Equal(Speed(111.19483768868857*KPH), to.SpeedFrom(from))
}

func TestSpeedFuncs(t *testing.T) {
	var assert = require.New(t)

	assert.True(Speed(MPH).FasterThan(KPH))
	assert.True(Speed(KPH).SlowerThan(MPH))
	assert.True(Speed(3600 * MPH).Equal(5280 * FeetPerSecond))
}

func TestParseDistance(t *testing.T) {
	var assert = require.New(t)

	assert.Zero(MustParseDistance(``))
	assert.Zero(MustParseDistance(`0m`))
	assert.Zero(MustParseDistance(`0 meter`))
	assert.Zero(MustParseDistance(`0 meters`))
	assert.Zero(MustParseDistance(`0meter`))
	assert.Zero(MustParseDistance(`0meters`))

	assert.EqualValues(1, MustParseDistance(`1m`))
	assert.EqualValues(1, MustParseDistance(`1 meter`))
	assert.EqualValues(1, MustParseDistance(`1 meters`))
	assert.EqualValues(1, MustParseDistance(`1meter`))
	assert.EqualValues(1, MustParseDistance(`1meters`))

	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625km`))
	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625 km`))
	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625 km.`))
	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625 kilometer`))
	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625 kilometers`))
	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625kilometer`))
	assert.EqualValues(Distance(3.141597625)*Kilometer, MustParseDistance(`3.141597625kilometers`))

	assert.EqualValues(Distance(26.2)*Mile, MustParseDistance(`26.2mi`))
	assert.EqualValues(Distance(26.2)*Mile, MustParseDistance(`26.2 miles`))

	assert.EqualValues(Distance(5280)*Foot, MustParseDistance(`5280ft`))
	assert.EqualValues(Distance(5280)*Foot, MustParseDistance(`5280 ft`))
	assert.EqualValues(Distance(5280)*Foot, MustParseDistance(`5280 ft.`))
	assert.EqualValues(Distance(5280)*Foot, MustParseDistance(`5280 feet`))
	assert.EqualValues(Distance(5280)*Foot, MustParseDistance(`5280feet`))

	assert.EqualValues(Distance(100)*Yard, MustParseDistance(`100yd.`))
	assert.EqualValues(Distance(100)*Yard, MustParseDistance(`100 yd.`))
	assert.EqualValues(Distance(100)*Yard, MustParseDistance(`100 yd`))
	assert.EqualValues(Distance(100)*Yard, MustParseDistance(`100yd`))
	assert.EqualValues(Distance(100)*Yard, MustParseDistance(`100 yards`))

	assert.EqualValues(Distance(300)*NauticalMile, MustParseDistance(`300nm`))
	assert.EqualValues(Distance(300)*NauticalMile, MustParseDistance(`300 nm.`))
	assert.EqualValues(Distance(300)*NauticalMile, MustParseDistance(`300 nm`))
	assert.EqualValues(Distance(300)*NauticalMile, MustParseDistance(`300 nautical miles`))
	assert.EqualValues(Distance(300)*NauticalMile, MustParseDistance(`300 nautical mile`))
}
