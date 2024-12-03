// Utilities for messing with time.
package timeutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/utils"
)

// Return the standard Golang reference time (2006-01-02T15:04:05.999999999Z07:00)
func ReferenceTime() time.Time {
	return utils.ReferenceTime
}

// Provides an API-compatible version of time.ParseDuration that accepts additional
// formats for parsing durations:
//
//	1y, 3year, 5years: Expands to (n*24*365) hours
//	1w, 3weeks, 5wks:  Expands to (n*24*7) hours
//	1d, 1day, 5days:   Expands to (n*24) hours
func ParseDuration(in string) (time.Duration, error) {
	return utils.ParseDuration(in)
}

// Formats the given duration in a colon-separated timer format in the form
// [HH:]MM:SS.
func FormatTimer(d time.Duration) string {
	h, m, s := DurationHMS(d)

	out := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	out = strings.TrimPrefix(out, `00:`)
	out = strings.TrimPrefix(out, `0`)
	return out
}

// Formats the given duration using the given format string.  The string follows
// the same formatting rules as described in the fmt package, and will receive
// three integer arguments: hours, minutes, and seconds.
func FormatTimerf(format string, d time.Duration) string {
	h, m, s := DurationHMS(d)

	out := fmt.Sprintf(format, h, m, s)
	return out
}

// Extracts the hours, minutes, and seconds from the given duration.
func DurationHMS(d time.Duration) (int, int, int) {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	return int(h), int(m), int(s)
}
