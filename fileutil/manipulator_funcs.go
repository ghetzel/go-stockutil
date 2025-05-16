package fileutil

import (
	"strings"
)

// Performs multiple sequential manipulations on an intercepted line of text from
// an io.Reader as its being read.
func ManipulateAll(fns ...ReadManipulatorFunc) ReadManipulatorFunc {
	return func(data []byte) ([]byte, error) {
		var err error

		for _, fn := range fns {
			data, err = fn(data)

			if err != nil {
				break
			}
		}

		return data, err
	}
}

// A ReadManipulatorFunc for replacing text in an io.Reader as its being read.
func ReplaceWith(find string, replace string, occurrences int) ReadManipulatorFunc {
	return func(data []byte) ([]byte, error) {
		var line = string(data)
		line = strings.Replace(line, find, replace, occurrences)
		return []byte(line), nil
	}
}

// A ReadManipulatorFunc for removing lines that only contain whitespace.
func RemoveBlankLines(data []byte) ([]byte, error) {
	if line := strings.TrimSpace(string(data)); len(line) == 0 {
		return nil, ErrSkipToken
	} else {
		return data, nil
	}
}

// A ReadManipulatorFunc for removing lines that contain the given string.
func RemoveLinesContaining(needle string) ReadManipulatorFunc {
	return func(data []byte) ([]byte, error) {
		if strings.Contains(string(data), needle) {
			return nil, nil
		} else {
			return data, nil
		}
	}
}

// A ReadManipulatorFunc for removing lines that have a given prefix.
func RemoveLinesWithPrefix(prefix string, trimSpace bool) ReadManipulatorFunc {
	return func(data []byte) ([]byte, error) {
		var line = string(data)

		if trimSpace {
			line = strings.TrimSpace(line)
		}

		if strings.HasPrefix(line, prefix) {
			return nil, nil
		} else {
			return data, nil
		}
	}
}

// A ReadManipulatorFunc for removing lines that have a given suffix.
func RemoveLinesWithSuffix(suffix string, trimSpace bool) ReadManipulatorFunc {
	return func(data []byte) ([]byte, error) {
		var line = string(data)

		if trimSpace {
			line = strings.TrimSpace(line)
		}

		if strings.HasSuffix(line, suffix) {
			return nil, nil
		} else {
			return data, nil
		}
	}
}

// A ReadManipulatorFunc for removing lines surrounded by a given prefix and suffix.
func RemoveLinesSurroundedBy(prefix string, suffix string, trimSpace bool) ReadManipulatorFunc {
	return func(data []byte) ([]byte, error) {
		var line = string(data)

		if trimSpace {
			line = strings.TrimSpace(line)
		}

		if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, suffix) {
			return nil, nil
		} else {
			return data, nil
		}
	}
}
