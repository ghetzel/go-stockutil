package log

import (
	"time"
)

type Timing struct {
	name    string
	started time.Time
	took    time.Duration
}

func (self *Timing) Reset() {
	self.took = 0
	self.started = time.Now()
}

func (self *Timing) Done() time.Duration {
	self.took = time.Since(self.started)
	Debugf("[${red+b}TIMING${reset}] %s took %v", self.name, self.took)
	return self.took
}

func (self *Timing) Then(name string) *Timing {
	self.Done()
	self.name = name
	self.Reset()
	return self
}

func Time(name string) *Timing {
	return &Timing{
		name:    name,
		started: time.Now(),
	}
}

func TimeFunc(name string, fn func()) *Timing {
	if fn == nil {
		panic("Cannot call log.TimeFunc with a nil function")
	}

	var tm = Time(name)
	fn()
	tm.Done()
	return tm
}
