package log

import (
	"fmt"
	"io"
	"regexp"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/mgutz/ansi"
)

var rxColorExpr = regexp.MustCompile(`(\$\{(?P<color>[^\}]+)\})`) // ${color}, ${color:mod1:mod2}
var TerminalEscapePrefix = `\[`
var TerminalEscapeSuffix = `\]`

func csprintf(termEscape bool, colorEnabled bool, format string, args ...any) string {
	var out = fmt.Sprintf(format, args...)

	for {
		if match := rxutil.Match(rxColorExpr, out); match != nil {
			var colorExpr = match.Group(`color`)
			var repl = ``

			// only replace with the actual ANSI escape sequences if we're at a tty
			// or if colors have been explicitly enabled, otherwise just remove the sequences
			if colorEnabled {
				repl = ansi.ColorCode(colorExpr)

				if termEscape {
					repl = stringutil.Wrap(repl, TerminalEscapePrefix, TerminalEscapeSuffix)
				}
			}

			out = match.ReplaceGroup(1, repl)
		} else {
			break
		}
	}

	return out
}

func CSprintf(format string, args ...any) string {
	return csprintf(false, true, format, args...)
}

func CFPrintf(w io.Writer, format string, args ...any) (int, error) {
	return fmt.Fprint(w, CSprintf(format, args...))
}

func CStripf(format string, args ...any) string {
	return csprintf(false, false, format, args...)
}

// Same as CSprintf, but wraps all replaced color sequences with terminal escape sequences
// as defined in TerminalEscapePrefix and TerminalEscapeSuffix
func TermSprintf(format string, args ...any) string {
	return csprintf(true, true, format, args...)
}
