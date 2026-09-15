package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	var assert = require.New(t)

	v, err := ParseDuration(``)
	assert.NoError(err)
	assert.Zero(v)

	for in, out := range map[string]time.Duration{
		`4h`:    time.Duration(4 * time.Hour),
		`4H`:    time.Duration(4 * time.Hour),
		`1d`:    time.Duration(24 * time.Hour),
		`1day`:  time.Duration(24 * time.Hour),
		`1days`: time.Duration(24 * time.Hour),
		`5 years 4 weeks 3 days 2 hours 1 minute`: time.Duration(44546*time.Hour) + time.Minute,
		`1d1h`:  time.Duration(25 * time.Hour),
		`1d 1h`: time.Duration(25 * time.Hour),
	} {
		v, err := ParseDuration(in)
		assert.NoError(err)
		assert.Equal(out, v, fmt.Sprintf("in=%v", in))
	}
}

func TestDurationHMS(t *testing.T) {
	var assert = require.New(t)

	h, m, s := DurationHMS(0)
	assert.Equal(0, h)
	assert.Equal(0, m)
	assert.Equal(0, s)

	h, m, s = DurationHMS(time.Second)
	assert.Equal(0, h)
	assert.Equal(0, m)
	assert.Equal(1, s)

	h, m, s = DurationHMS(time.Minute)
	assert.Equal(0, h)
	assert.Equal(1, m)
	assert.Equal(0, s)

	h, m, s = DurationHMS(time.Hour)
	assert.Equal(1, h)
	assert.Equal(0, m)
	assert.Equal(0, s)

	h, m, s = DurationHMS(59 * time.Second)
	assert.Equal(0, h)
	assert.Equal(0, m)
	assert.Equal(59, s)

	h, m, s = DurationHMS(119 * time.Second)
	assert.Equal(0, h)
	assert.Equal(1, m)
	assert.Equal(59, s)

	h, m, s = DurationHMS(23*time.Hour + 4*time.Minute + 13*time.Second)
	assert.Equal(23, h)
	assert.Equal(4, m)
	assert.Equal(13, s)
}

func TestFormatTimer(t *testing.T) {
	var assert = require.New(t)

	var out = FormatTimer(0)
	assert.Equal(`0:00`, out)

	out = FormatTimer(time.Second)
	assert.Equal(`0:01`, out)

	out = FormatTimer(time.Minute)
	assert.Equal(`1:00`, out)

	out = FormatTimer(time.Hour)
	assert.Equal(`1:00:00`, out)

	out = FormatTimer(59 * time.Second)
	assert.Equal(`0:59`, out)

	out = FormatTimer(119 * time.Second)
	assert.Equal(`1:59`, out)

	out = FormatTimer(23*time.Hour + 4*time.Minute + 13*time.Second)
	assert.Equal(`23:04:13`, out)
}

func ExampleFormatTimer_zeroValue() {
	fmt.Print(FormatTimer(0))
	// Output: 0:00
}

func ExampleFormatTimer_oneSecond() {
	fmt.Print(FormatTimer(time.Second))
	// Output: 0:01
}

func ExampleFormatTimer_oneMinute() {
	fmt.Print(FormatTimer(time.Minute))
	// Output: 1:00
}

func ExampleFormatTimer_oneHour() {
	fmt.Print(FormatTimer(time.Hour))
	// Output: 1:00:00
}

func ExampleFormatTimer_underOneMinute() {
	fmt.Print(FormatTimer(59 * time.Second))
	// Output: 0:59
}

func ExampleFormatTimer_overOneMinute() {
	fmt.Print(FormatTimer(119 * time.Second))
	// Output: 1:59
}

func ExampleFormatTimer_complete() {
	fmt.Print(FormatTimer(23*time.Hour + 4*time.Minute + 13*time.Second))
	// Output: 23:04:13
}

func tm(tm string) time.Time {
	var t, _ = time.Parse(time.RFC3339, tm)
	return t
}

func TestRoundFloor(t *testing.T) {
	var assert = require.New(t)

	type refRoundWanted struct {
		Ref   time.Time
		Round time.Duration
		Want  time.Time
		NEq   bool
	}

	for _, rrw := range []refRoundWanted{
		{tm(`2014-08-17T00:17:34Z`), time.Minute, tm(`2014-08-17T00:17:00Z`), false},
		{tm(`2014-08-17T00:17:34Z`), time.Minute, tm(`2014-08-17T00:18:00Z`), true},
		{tm(`2014-08-17T00:17:34Z`), 5 * time.Minute, tm(`2014-08-17T00:15:00Z`), false},
		{tm(`2014-08-17T00:17:34Z`), 5 * time.Minute, tm(`2014-08-17T00:20:00Z`), true},
		{tm(`2014-08-17T00:17:34Z`), 30 * time.Second, tm(`2014-08-17T00:17:30Z`), false},
		{tm(`2014-08-17T00:17:34Z`), 30 * time.Second, tm(`2014-08-17T00:17:20Z`), true},
		{tm(`2014-12-31T23:17:34Z`), time.Hour, tm(`2014-12-31T23:00:00Z`), false},
		{tm(`2014-12-31T23:17:34Z`), 24 * time.Hour, tm(`2014-12-31T00:00:00Z`), false},
		{tm(`2014-12-31T19:17:34-04:00`), 24 * time.Hour, tm(`2014-12-30T20:00:00-04:00`), false},
	} {
		if !rrw.NEq {
			assert.Equal(rrw.Want, RoundFloor(rrw.Ref, rrw.Round))
		} else {
			assert.NotEqual(rrw.Want, RoundFloor(rrw.Ref, rrw.Round))
		}
	}
}

func TestRoundCeil(t *testing.T) {
	var assert = require.New(t)

	type refRoundWanted struct {
		Ref   time.Time
		Round time.Duration
		Want  time.Time
		NEq   bool
	}

	for _, rrw := range []refRoundWanted{
		{tm(`2014-08-17T00:17:34Z`), time.Minute, tm(`2014-08-17T00:18:00Z`), false},
		{tm(`2014-08-17T00:17:34Z`), time.Minute, tm(`2014-08-17T00:20:00Z`), true},
		{tm(`2014-08-17T00:17:34Z`), 5 * time.Minute, tm(`2014-08-17T00:20:00Z`), false},
		{tm(`2014-08-17T00:17:34Z`), 5 * time.Minute, tm(`2014-08-17T00:18:00Z`), true},
		{tm(`2014-08-17T00:17:34Z`), 10 * time.Second, tm(`2014-08-17T00:17:40Z`), false},
		{tm(`2014-08-17T00:17:34Z`), 10 * time.Second, tm(`2014-08-17T00:17:20Z`), true},
		{tm(`2014-12-31T23:17:34Z`), time.Hour, tm(`2015-01-01T00:00:00Z`), false},
		{tm(`2014-12-31T23:17:34Z`), 24 * time.Hour, tm(`2015-01-01T00:00:00Z`), false},
		{tm(`2014-12-31T19:17:34-04:00`), 24 * time.Hour, tm(`2014-12-31T20:00:00-04:00`), false},
		{tm(`2014-12-31T23:17:34-04:00`), 24 * time.Hour, tm(`2015-01-01T20:00:00-04:00`), false},
	} {
		if !rrw.NEq {
			assert.Equal(rrw.Want, RoundCeil(rrw.Ref, rrw.Round))
		} else {
			assert.NotEqual(rrw.Want, RoundCeil(rrw.Ref, rrw.Round))
		}
	}
}
