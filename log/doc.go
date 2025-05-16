/*
Package log package provides convenient and flexible utilities for logging messages.

# Overview

Logging is done by calling functions in this package that correspond to the
severity of the log message being output.  At the package level, a minimum
severity can be set.  Messages less severe than this minimum logging level will
not be output.

# Color Expressions

In addition to the standard printf-style formatting options (as defined in the
standard fmt package), this package supports inline expressions that control the
output of ANSI terminal escape sequences.  These expressions allow for a simple
mechanism to colorize log output, as well as applying graphical effects like
bold, underline, and blinking text (for terminals that support it).

By default, color expressions will only be honored if os.Stdin is attached to a
pseudoterminal.  This is the case when the program is run on the command line
and is not piped or redirected to another file.  This default ensures that the
colors are visible only in a visual context, but do not corrupt files or
pipelines with ANSI escape sequences.  Color sequences can be explicitly enabled
or disabled by setting the EnableColorExpressions package variable.

Using color expressions in format strings is done by wrapping the expression in
${expr}.  The general format for color expressions is:

	foregroundColor[+attributes[:backgroundColor[+attributes]]]

Colors (foreground and background):

	black
	red
	green
	yellow
	blue
	magenta
	cyan
	white
	[0-255]: numeric 8-bit color (for 256 color terminals)
	reset: Reset all color and graphics attributes to their defaults

Foreground Attributes:

	b: bold text
	B: blinking text
	h: high-intensity (bright text)
	i: inverted/reverse colors
	s: strikethrough
	u: underline

Background Attributes:

	h: high-intensity (bright text)

# Examples

Below are some examples showing various formatting options for logs.

	log.Info("Hello, world!")
	// [11:22:33 0001] INFO Hello, world!

	log.Warningf("The %q operation could not be completed.", "add")
	// [11:22:33 0002] WARN The "add" operation could not be completed.

	log.Errorf("There was an ${red}error${reset} opening file ${blue+b:white}%s${reset}", filename)
	// [11:22:33 0003] ERRO There was an error opening file /tmp/file.txt
	//                                   ^^^^^              ^^^^^^^^^^^^^
	//                                   red text           blue text on white background

# Log Interception

It is sometimes useful to be able to act on logs as they are emitted, especially in cases where this
package is used in other projects that are imported.  Log Interceptors are called before each log
line is emitted.  The LogInterceptFunc is called with the level the message was emitted with, the
message itself as a string, and a stack trace struct that defines exactly where the log was emitted
from.

	// print a complete stack trace before every debug-level message that is encountered
	log.AddLogIntercept(func(level log.Level, line string, stack log.StackItems){
		if level == log.DEBUG {
			for _, item := range stack {
				fmt.Println(item.String())
			}
		}
	})

# Writable Logger

The WritableLogger implements the io.Writer interface, acting as a bridge between byte streams from
various sources and the log package.  This is frequently useful in situations like parsing the
output of other programs.  A WritableLogger accepts a custom LogParseFunc that allows individual
lines being written to the WritableLogger to be parsed, rewritten, and given a log severity level.

	import (
		"os/exec"
		"github.com/ghetzel/go-stockutil/log"
	)

	var wr = log.NewWritableLogger(log.INFO, `ls: `)

	wr.SetParserFunc(func(line string) (log.Level, string) {
		if strings.Contains(line, `root`) {
			// root-owned files show up as errors
			return log.ERROR, line
		} else if strings.Contains(line, os.Getenv(`USER`)) {
			// current user files are notices
			return log.NOTICE, line
		} else {
			// all other lines are not logged at all
			return log.DEBUG, ``
		}
	})

	var ls = exec.Command(`ls`, `-l`)
	ls.Stdout = wr
	ls.Run()
*/
package log
