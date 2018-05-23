package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
	logging "github.com/op/go-logging"
)

const defaultFormat = `[%{time:15:04:05} %{id:04d}] %{color:bold}%{level:.4s}%{color:reset} %{message}`

type Logger struct {
	module    string
	dest      io.Writer
	format    string
	sequence  uint64
	level     Level
	formatter logging.Formatter
}

func NewLogger(module string) *Logger {
	return NewLoggerWriter(module, os.Stderr)
}

func NewLoggerWriter(module string, w io.Writer) *Logger {
	logger := &Logger{
		module: module,
		level:  DEBUG,
	}

	if err := logger.SetFormat(defaultFormat); err != nil {
		panic(err.Error())
	}

	return logger
}

func (self *Logger) SetFormat(format string) error {
	if f, err := logging.NewStringFormatter(format); err == nil {
		self.formatter = f
		return nil
	} else {
		return err
	}
}

func (self *Logger) SetLevel(level Level) {
	self.level = level
}

func (self *Logger) ShouldLog(level Level) bool {
	if level <= self.level {
		return true
	} else {
		return false
	}
}

func (self *Logger) Logf(level Level, format string, args ...interface{}) {
	if self.ShouldLog(level) {
		message := fmt.Sprintf(format, args...)

		self.formatter.Format(0, &logging.Record{
			ID:     atomic.AddUint64(&self.sequence, 1),
			Time:   time.Now(),
			Module: self.module,
			Level:  logging.Level(level),
			Args:   []interface{}{message},
		}, self.dest)
	}
}

func (self *Logger) Log(level Level, args ...interface{}) {
	Logf(level, "%v", args...)
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
