package log

import (
	"strings"
)

type LogParseFunc func(line string) (Level, string)

type WritableLogger struct {
	level   Level
	prefix  string
	levelFn LogParseFunc
}

func NewWritableLogger(level Level, prefix ...string) *WritableLogger {
	return &WritableLogger{
		level:  level,
		prefix: strings.Join(prefix, ``),
	}
}

func (self *WritableLogger) SetParserFunc(fn LogParseFunc) *WritableLogger {
	self.levelFn = fn
	return self
}

func (self *WritableLogger) Write(p []byte) (int, error) {
	initLogging()

	var lvl = self.level
	var line = string(p)

	if self.levelFn != nil {
		if newLevel, rewritten := self.levelFn(line); rewritten != `` {
			lvl = newLevel
			line = rewritten
		} else {
			return len(p), nil
		}
	}

	Logf(lvl, "%v%v", self.prefix, line)
	return len(p), nil
}
