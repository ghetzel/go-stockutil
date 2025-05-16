package rxutil

// Splits a given string using the given regexp string or *regexp.Regexp value into at most n parts.
func SplitN(pattern any, s string, n int) []string {
	if s == `` {
		return make([]string, 0)
	}

	if match := Match(pattern, s); match != nil {
		return match.rx.Split(s, n)
	} else {
		return []string{s}
	}
}

// Splits a given string using the given regexp string or *regexp.Regexp value into zero or more parts.
func Split(pattern any, s string) []string {
	return SplitN(pattern, s, -1)
}
