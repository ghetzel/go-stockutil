package log

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/mgutz/ansi"
)

var rxColorExpr = regexp.MustCompile(`(\$\{(?P<color>[^\}]+)\})`) // ${color}, ${color:mod1:mod2}
var TerminalEscapePrefix = `\[`
var TerminalEscapeSuffix = `\]`

func csprintf(termEscape bool, colorEnabled bool, format string, literal string, args ...any) string {
	var out string

	if literal != `` {
		out = literal
	} else if format != `` {
		out = fmt.Sprintf(format, args...)
	} else {
		out = fmt.Sprint(args...)
	}

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

// A version of fmt.Sprint that supports terminal color expressions.
func CSprint(messages ...any) string {
	return csprintf(false, true, ``, ``, messages...)
}

// A version of fmt.Sprintf that supports terminal color expressions.
func CSprintf(format string, args ...any) string {
	return csprintf(false, true, format, ``, args...)
}

// A version of fmt.Fprint that supports terminal color expressions.
func CFprint(w io.Writer, messages ...any) (int, error) {
	return fmt.Fprint(w, CSprint(messages...))
}

// A version of fmt.Fprintf that supports terminal color expressions.
func CFprintf(w io.Writer, format string, args ...any) (int, error) {
	return fmt.Fprint(w, CSprintf(format, args...))
}

// Alias of log.CFprintf
func CFPrintf(w io.Writer, format string, args ...any) (int, error) {
	return CFprintf(w, format, args...)
}

// A version of fmt.Print that supports terminal color expressions.
func Cprint(messages ...any) (int, error) {
	return CFprint(os.Stdout, messages...)
}

// A version of fmt.Printf that supports terminal color expressions.
func Cprintf(format string, args ...any) (int, error) {
	return CFprintf(os.Stdout, format, args...)
}

// Return a string with all color expressions removed.
func CStripf(format string, args ...any) string {
	return csprintf(false, false, format, ``, args...)
}

// Same as CSprintf, but wraps all replaced color sequences with terminal escape sequences
// as defined in TerminalEscapePrefix and TerminalEscapeSuffix
func TermSprintf(format string, args ...any) string {
	return csprintf(true, true, format, ``, args...)
}
