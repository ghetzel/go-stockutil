package log

import (
	"fmt"
	"io"
	"os"

	"github.com/ghetzel/go-stockutil/typeutil"
	logging "github.com/op/go-logging"
)

type Logger struct {
	module    string
	backend   *logging.LogBackend
	formatted logging.Backend
	leveled   logging.LeveledBackend
	logger    *logging.Logger
}

func NewLogger(module string) *Logger {
	return NewLoggerWriter(module, os.Stderr)
}

func NewLoggerWriter(module string, w io.Writer) *Logger {
	logger := &Logger{
		module:  module,
		backend: logging.NewLogBackend(w, ``, 0),
	}

	logger.formatted = logging.NewBackendFormatter(logger.backend, logging.MustStringFormatter(
		fmt.Sprintf(
			`[%%{time:15:04:05} %%{id:04d}] %%{color:bold}%%{level:.4s}%%{color:reset} %%{message}`,
		),
	))

	logger.leveled = logging.AddModuleLevel(logger.formatted)
	logger.logger = logging.MustGetLogger(module)

	return logger
}

func (self *Logger) SetLevel(level Level) error {
	if lvl, err := logging.LogLevel(level.String()); err == nil {
		self.leveled.SetLevel(lvl, self.module)
		return nil
	} else {
		return fmt.Errorf("[INVALID LEVEL %v] ", level)
	}
}

func (self *Logger) Logf(level Level, format string, args ...interface{}) {
	switch level {
	case PANIC:
		self.logger.Panicf(format, args...)
	case FATAL:
		self.logger.Fatalf(format, args...)
	case CRITICAL:
		self.logger.Criticalf(format, args...)
	case ERROR:
		self.logger.Errorf(format, args...)
	case WARNING:
		self.logger.Warningf(format, args...)
	case NOTICE:
		self.logger.Noticef(format, args...)
	case INFO:
		self.logger.Infof(format, args...)
	default:
		self.logger.Debugf(format, args...)
	}
}

func (self *Logger) Log(level Level, args ...interface{}) {
	switch level {
	case PANIC:
		self.logger.Panic(args...)
	case FATAL:
		self.logger.Fatal(args...)
	case CRITICAL:
		self.logger.Critical(args...)
	case ERROR:
		self.logger.Error(args...)
	case WARNING:
		self.logger.Warning(args...)
	case NOTICE:
		self.logger.Notice(args...)
	case INFO:
		self.logger.Info(args...)
	default:
		self.logger.Debug(args...)
	}
}

func (self *Logger) Critical(args ...interface{}) {
	self.Log(CRITICAL, args...)
}

func (self *Logger) Criticalf(format string, args ...interface{}) {
	self.Logf(CRITICAL, format, args...)
}

func (self *Logger) Debug(args ...interface{}) {
	self.Log(DEBUG, args...)
}

func (self *Logger) Debugf(format string, args ...interface{}) {
	self.Logf(DEBUG, format, args...)
}

func (self *Logger) Error(args ...interface{}) {
	self.Log(ERROR, args...)
}

func (self *Logger) Errorf(format string, args ...interface{}) {
	self.Logf(ERROR, format, args...)
}

func (self *Logger) Fatal(args ...interface{}) {
	self.Log(FATAL, args...)
}

func (self *Logger) Fatalf(format string, args ...interface{}) {
	self.Logf(FATAL, format, args...)
}

func (self *Logger) Info(args ...interface{}) {
	self.Log(INFO, args...)
}

func (self *Logger) Infof(format string, args ...interface{}) {
	self.Logf(INFO, format, args...)
}

func (self *Logger) Notice(args ...interface{}) {
	self.Log(NOTICE, args...)
}

func (self *Logger) Noticef(format string, args ...interface{}) {
	self.Logf(NOTICE, format, args...)
}

func (self *Logger) Panic(args ...interface{}) {
	self.Log(PANIC, args...)
}

func (self *Logger) Panicf(format string, args ...interface{}) {
	self.Logf(PANIC, format, args...)
}

func (self *Logger) Warning(args ...interface{}) {
	self.Log(WARNING, args...)
}

func (self *Logger) Warningf(format string, args ...interface{}) {
	self.Logf(WARNING, format, args...)
}

func (self *Logger) Dump(args ...interface{}) {
	for _, arg := range args {
		self.Log(DEBUG, typeutil.Dump(arg))
	}
}

func (self *Logger) Dumpf(format string, args ...interface{}) {
	for _, arg := range args {
		self.Logf(DEBUG, format, typeutil.Dump(arg))
	}
}
