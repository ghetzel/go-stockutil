// Utilities to make working with regular expressions easier.
package rxutil

import (
	"fmt"
	"regexp"
)

type MatchResult struct {
	rx                        *regexp.Regexp
	source                    string
	names                     []string
	leftmost_submatches       []string
	leftmost_submatch_indices []int
	all_submatches            [][]string
}

func IsMatchString(pattern any, source string) bool {
	return (Match(pattern, source) != nil)
}

// Returns a MatchResult object representing the leftmost match of pattern against
// source, or nil if no matches were found.  Pattern can be a string or a
// previously-compiled *regexp.Regexp.
func Match(pattern any, source string) *MatchResult {
	var rx *regexp.Regexp

	if r, ok := pattern.(*regexp.Regexp); ok {
		rx = r
	} else if r, ok := pattern.(regexp.Regexp); ok {
		rx = &r
	} else {
		rx = regexp.MustCompile(fmt.Sprintf("%v", pattern))
	}

	if rx.MatchString(source) {
		return &MatchResult{
			rx:                        rx,
			source:                    source,
			names:                     rx.SubexpNames(),
			leftmost_submatches:       rx.FindStringSubmatch(source),
			leftmost_submatch_indices: rx.FindStringSubmatchIndex(source),
			all_submatches:            rx.FindAllStringSubmatch(source, -1),
		}
	}

	return nil
}

func (self *MatchResult) groupI(nameOrIndex any) (string, int) {
	for i, name := range self.names {
		switch nameOrIndex.(type) {
		case string:
			if i > 0 && name == nameOrIndex.(string) {
				return self.leftmost_submatches[i], i
			}
		case int:
			if i == nameOrIndex.(int) {
				return self.leftmost_submatches[i], i
			}
		}
	}

	return ``, -1
}

// Return the value of the numbered capture group (if given an int), or the
// named capture group (if given a string).  Returns an empty string if the
// given group name or index does not exist.
func (self *MatchResult) Group(nameOrIndex any) string {
	if match, i := self.groupI(nameOrIndex); i >= 0 {
		return match
	} else {
		return ``
	}
}

// Return a copy of source string with the given numbered or named group replaced
// with repl.
func (self *MatchResult) ReplaceGroup(nameOrIndex any, repl string) string {
	if _, i := self.groupI(nameOrIndex); i >= 0 {
		if i == 0 {
			return repl
		}

		if (i*2 + 1) < len(self.leftmost_submatch_indices) {
			var startIndex = self.leftmost_submatch_indices[(i * 2)]
			var endIndex = self.leftmost_submatch_indices[(i*2)+1]

			var out = self.source[0:startIndex]
			out += repl
			out += self.source[endIndex:]

			return out
		}
	}

	return self.source
}

// Returns a map of all named capture matches, keyed on capture group name.
func (self *MatchResult) NamedCaptures() map[string]string {
	var captures = make(map[string]string)

	for i, name := range self.names {
		if i > 0 && name != `` {
			captures[name] = self.leftmost_submatches[i]
		}
	}

	return captures
}

// Return a slice of all capture groups.
func (self *MatchResult) Captures() []string {
	return self.leftmost_submatches
}

// Returns all captures from all matches appended together.  The full match string
// from match is omitted, so only the actual values appearing within capture groups
// are returned.
func (self *MatchResult) AllCaptures() []string {
	var all = make([]string, 0)

	for _, subcaptures := range self.all_submatches {
		if len(subcaptures) > 1 {
			all = append(all, subcaptures[1:]...)
		}
	}

	return all
}
