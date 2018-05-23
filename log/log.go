// Standard logging package, batteries included.
package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghetzel/go-stockutil/typeutil"
	multierror "github.com/hashicorp/go-multierror"
)

var defaultLogger *Logger
var ModuleName = ``

type Loggable interface {
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

type LogFunc func(args ...interface{})
type FormattedLogFunc func(format string, args ...interface{})

// The LOGLEVEL environment variable has final say over the effective log level
// for all users of this package.
var LogLevel Level = func() Level {
	if v := os.Getenv(`LOGLEVEL`); v != `` {
		return GetLevel(v)
	} else {
		return INFO
	}
}()

func initLogging() {
	if defaultLogger == nil {
		defaultLogger = NewLogger(``)
		SetLevel(LogLevel)
	}
}

func Debugging() bool {
	return (LogLevel == DEBUG)
}

func DefaultLogger() Loggable {
	initLogging()
	return defaultLogger
}

func SetLevelString(level string, modules ...string) {
	SetLevel(GetLevel(level), modules...)
}

func SetLevel(level Level, modules ...string) {
	initLogging()
	defaultLogger.SetLevel(level)
	LogLevel = level
}

func Logf(level Level, format string, args ...interface{}) {
	initLogging()
	defaultLogger.Logf(level, format, args...)
}

func Log(level Level, args ...interface{}) {
	initLogging()
	defaultLogger.Log(level, args...)
}

func Critical(args ...interface{}) {
	Log(CRITICAL, args...)
}

func Criticalf(format string, args ...interface{}) {
	Logf(CRITICAL, format, args...)
}

func Debug(args ...interface{}) {
	Log(DEBUG, args...)
}

func Debugf(format string, args ...interface{}) {
	Logf(DEBUG, format, args...)
}

func Dump(args ...interface{}) {
	for _, arg := range args {
		Log(DEBUG, typeutil.Dump(arg))
	}
}

func Dumpf(format string, args ...interface{}) {
	for _, arg := range args {
		Logf(DEBUG, format, typeutil.Dump(arg))
	}
}

func Error(args ...interface{}) {
	Log(ERROR, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(ERROR, format, args...)
}

func Fatal(args ...interface{}) {
	Log(FATAL, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logf(FATAL, format, args...)
}

func Info(args ...interface{}) {
	Log(INFO, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(INFO, format, args...)
}

func Notice(args ...interface{}) {
	Log(NOTICE, args...)
}

func Noticef(format string, args ...interface{}) {
	Logf(NOTICE, format, args...)
}

func Panic(args ...interface{}) {
	Log(PANIC, args...)
}

func Panicf(format string, args ...interface{}) {
	Logf(PANIC, format, args...)
}

func Warning(args ...interface{}) {
	Log(WARNING, args...)
}

func Warningf(format string, args ...interface{}) {
	Logf(WARNING, format, args...)
}

func Confirm(prompt string) bool {
	return Confirmf(prompt)
}

func Confirmf(format string, args ...interface{}) bool {
	var response string

	fmt.Printf(format, args...)

	if _, err := fmt.Scanln(&response); err != nil {
		panic(err.Error())
	}

	for _, okay := range []string{
		`y`,
		`yes`,
	} {
		if strings.ToLower(okay) == strings.ToLower(response) {
			return true
		}
	}

	return false
}

// Appends on error to another, allowing for operations that return multiple errors
// to remain compatible within a single-valued context.
func AppendError(base error, err error) error {
	if err == nil {
		return base
	} else {
		return multierror.Append(base, err)
	}
}
