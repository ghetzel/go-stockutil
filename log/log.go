package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	multierror "github.com/hashicorp/go-multierror"
	isatty "github.com/mattn/go-isatty"
	"github.com/op/go-logging"
)

var EnableColorExpressions = func() bool {
	if forceColor := os.Getenv(`FORCE_COLOR`); forceColor != `` {
		return typeutil.Bool(forceColor)
	} else {
		return isatty.IsTerminal(os.Stdout.Fd())
	}
}()

var DefaultInterceptStackDepth int = 5
var SynchronousIntercepts = false

var backend *logging.LogBackend
var formatted logging.Backend
var leveled logging.LeveledBackend
var intercepts sync.Map
var filters sync.Map
var defaultLogger *logging.Logger
var ModuleName = ``

type LogFunc func(args ...any)
type FormattedLogFunc func(format string, args ...any)
type LogInterceptFunc func(level Level, line string, stack StackItems)
type LogFilterFunc func(level Level, line string, stack StackItems) bool

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
		backend = logging.NewLogBackend(os.Stderr, ``, 0)

		formatted = logging.NewBackendFormatter(backend, logging.MustStringFormatter(
			fmt.Sprintf(
				`[%%{time:15:04:05} %%{id:04d}] %%{color:bold}%%{level:.4s}%%{color:reset} %%{message}`,
			),
		))

		leveled = logging.AddModuleLevel(formatted)
		logging.SetBackend(leveled)

		defaultLogger = logging.MustGetLogger(ModuleName)
		SetLevel(LogLevel)
	}
}

// Append a function to be called (asynchronously in its own goroutine, or
// synchronously if SynchronousIntercepts is true) for every line logged.
// Returns a UUID that can be later used to deregister the intercept function.
func AddLogIntercept(fn LogInterceptFunc) string {
	var id = stringutil.UUID().String()
	intercepts.Store(id, fn)
	return id
}

// Remove the previously-added log intercept function.
func RemoveLogIntercept(id string) {
	intercepts.Delete(id)
}

// Append a function to be called for every line logged.  If the function
// returns false, the line will be dropped (will not be logged or intercepted).
func AddLogFilter(fn LogFilterFunc) string {
	var id = stringutil.UUID().String()
	filters.Store(id, fn)
	return id
}

// Remove the previously-added log filter function.
func RemoveLogFilter(id string) {
	filters.Delete(id)
}

// Takes a list of package names ("a/b/c") or patterns ("a/b/*") whose log output
// will be rejected.  If the package name or patterns matches any function in the
// call stack, the associated message will be rejected.
func FilterPackages(patterns ...string) string {
	return AddLogFilter(func(level Level, line string, stack StackItems) bool {
		for _, item := range stack {
			for _, pkg := range patterns {
				if item.InPackage(pkg) {
					return false
				}
			}
		}

		return true
	})
}

func Debugging() bool {
	return (LogLevel == DEBUG)
}

func VeryDebugging(features ...string) bool {
	if Debugging() {
		var envFeatures = strings.Split(os.Getenv(`DEBUG`), `,`)

		for _, feature := range features {
			for _, ef := range envFeatures {
				if typeutil.Bool(ef) {
					return true
				} else if strings.ToLower(feature) == strings.ToLower(ef) {
					return true
				}
			}
		}
	}

	return false
}

func Logger() *logging.Logger {
	initLogging()
	return defaultLogger
}

// Set the destination Writer where logs will henceforth be written.
func SetOutput(w io.Writer) {
	initLogging()
	backend.Logger.SetOutput(w)
}

func SetLevelString(level string, modules ...string) {
	SetLevel(GetLevel(level), modules...)
}

func SetLevel(level Level, modules ...string) {
	initLogging()

	if lvl, err := logging.LogLevel(level.String()); err == nil {
		if len(modules) == 0 {
			leveled.SetLevel(lvl, ``)
		} else {
			for _, module := range modules {
				leveled.SetLevel(lvl, module)
			}
		}
	} else {
		fmt.Printf("[INVALID LEVEL %v] ", level)
	}
}

func Logf(level Level, format string, args ...any) {
	initLogging()

	var line = fmt.Sprintf(format, args...)
	var stack = StackTrace(DefaultInterceptStackDepth)

	if shouldSkip(level, line, stack) {
		return
	}

	callIntercepts(level, line, stack)

	// only replace with the actual ANSI escape sequences if we're at a tty
	// or if colors have been explicitly enabled, otherwise just remove the sequences
	if EnableColorExpressions {
		log(level, CSprintf(format, args...))
	} else {
		log(level, CStripf(format, args...))
	}
}

func Log(level Level, args ...any) {
	// handle this special case where we are handling a fatal-level unformatted nil value,
	// in which case we don't actually want to end the program.
	//
	// NOTE: this obviates the need for FatalfIf, making Fatal() behave in the same way
	//
	if level == FATAL && len(args) == 1 && args[0] == nil {
		return
	}

	initLogging()

	var line = strings.Join(sliceutil.Stringify(args), ` `)
	var stack = StackTrace(DefaultInterceptStackDepth)

	if shouldSkip(level, line, stack) {
		return
	}

	callIntercepts(level, line, stack)
	log(level, args...)
}

func log(level Level, args ...any) {
	switch level {
	case PANIC:
		defaultLogger.Panic(args...)
	case FATAL:
		defaultLogger.Fatal(args...)
	case CRITICAL:
		defaultLogger.Critical(args...)
	case ERROR:
		defaultLogger.Error(args...)
	case WARNING:
		defaultLogger.Warning(args...)
	case NOTICE:
		defaultLogger.Notice(args...)
	case INFO:
		defaultLogger.Info(args...)
	default:
		defaultLogger.Debug(args...)
	}
}

func Critical(args ...any) {
	Log(CRITICAL, args...)
}

func Criticalf(format string, args ...any) {
	Logf(CRITICAL, format, args...)
}

func Debug(args ...any) {
	Log(DEBUG, args...)
}

func Debugf(format string, args ...any) {
	Logf(DEBUG, format, args...)
}

// Pretty-print the given arguments to the log at debug-level.
func Dump(args ...any) {
	for _, arg := range args {
		Log(LogLevel, typeutil.Dump(arg))
	}
}

// Same as Dump, but accepts a format string.
func Dumpf(format string, args ...any) {
	for _, arg := range args {
		Logf(LogLevel, format, typeutil.Dump(arg))
	}
}

// Marshal the arguments as indented JSON and log them at debug-level.
func DumpJSON(args ...any) {
	for _, arg := range args {
		if data, err := json.MarshalIndent(arg, ``, `  `); err == nil {
			Log(LogLevel, string(data))
		} else {
			Logf(LogLevel, "DumpJSON: %v", err)
		}
	}
}

func Error(args ...any) {
	Log(ERROR, args...)
}

func Errorf(format string, args ...any) {
	Logf(ERROR, format, args...)
}

func Fatal(args ...any) {
	Log(FATAL, args...)
}

func Fatalf(format string, args ...any) {
	Logf(FATAL, format, args...)
}

func Info(args ...any) {
	Log(INFO, args...)
}

func Infof(format string, args ...any) {
	Logf(INFO, format, args...)
}

func Notice(args ...any) {
	Log(NOTICE, args...)
}

func Noticef(format string, args ...any) {
	Logf(NOTICE, format, args...)
}

func Panic(args ...any) {
	Log(PANIC, args...)
}

func Panicf(format string, args ...any) {
	Logf(PANIC, format, args...)
}

func Warning(args ...any) {
	Log(WARNING, args...)
}

func Warningf(format string, args ...any) {
	Logf(WARNING, format, args...)
}

func Confirm(prompt string) bool {
	return Confirmf("%v", prompt)
}

// Present a confirmation prompt. The function returns true if the user interactively responds
// with "yes" or "y". Otherwise the function returns false.
func Confirmf(format string, args ...any) bool {
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

// Appends one error to another, allowing for operations that return multiple errors
// to remain compatible within a single-valued context.
func AppendError(base error, err error) error {
	if err == nil {
		return base
	} else {
		return multierror.Append(base, err)
	}
}

// Invoke Fatal() if the given error is not nil.
func FatalIf(err error) {
	if err != nil {
		Fatal(err)
	}
}

// Invoke Fatalf() if the given error is not nil.
func FatalfIf(format string, err error) {
	if err != nil {
		Fatalf(format, err)
	}
}

// call all registered intercept functions using the given arguments.
func callIntercepts(level Level, line string, stack StackItems) {
	intercepts.Range(func(_ any, value any) bool {
		if fn, ok := value.(LogInterceptFunc); ok {
			// for levels CRITICAL and worse, call intercepts synchronously in case we're
			// panicking and about to tear crap down.  Since these intercepts should run BEFORE
			// the log line is emitted, this should ensure the intercept definitely runs before
			// any of that goes down.
			if level <= CRITICAL || SynchronousIntercepts {
				fn(level, line, stack)
			} else {
				go fn(level, line, stack)
			}
		}

		return true
	})
}

func shouldSkip(level Level, line string, stack StackItems) (skip bool) {
	filters.Range(func(_ any, value any) bool {
		if fn, ok := value.(LogFilterFunc); ok {
			if !fn(level, line, stack) {
				skip = true
				return false
			}
		}

		return true
	})

	return
}
