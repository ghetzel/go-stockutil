package utils

import (
	"io"
	"time"
)

// A TimedReadCloser wraps an io.ReadCloser, keeping track of how long actually
// reading from and closing the ReadCloser took, as well as how many bytes were read.
//
// Measurement starts from the first call to Read(), and ends when Close() is called.
type TimedReadCloser struct {
	rc        io.ReadCloser
	bytesRead int64
	startedAt time.Time
	lastRead  time.Time
	readTook  time.Duration
	closeTook time.Duration
	took      time.Duration
}

func NewTimedReadCloser(rc io.ReadCloser) *TimedReadCloser {
	return &TimedReadCloser{
		rc: rc,
	}
}

func (self *TimedReadCloser) Read(b []byte) (int, error) {
	if self.startedAt.IsZero() {
		self.startedAt = time.Now()
	}

	n, err := self.rc.Read(b)
	self.bytesRead += int64(n)
	self.readTook = time.Since(self.startedAt)
	self.lastRead = time.Now()

	return n, err
}

func (self *TimedReadCloser) Close() (err error) {
	var closeStart = time.Now()
	self.readTook = time.Since(self.startedAt)

	err = self.rc.Close()

	self.closeTook = time.Since(closeStart)
	self.took = time.Since(self.startedAt)

	return
}

// Reset the internal counters to zero.  Useful after calling SetReadCloser().
func (self *TimedReadCloser) Reset() {
	self.startedAt = time.Time{}
	self.readTook = 0
	self.closeTook = 0
	self.took = 0
	self.bytesRead = 0
}

// Set the underlying io.ReadCloser.  Does not call Reset(), so multiple ReadClosers
// can be cumulitively tracked.
func (self *TimedReadCloser) SetReadCloser(rc io.ReadCloser) {
	self.rc = rc
}

// Return the time of the first call to Read().
func (self *TimedReadCloser) StartedAt() time.Time {
	return self.startedAt
}

// Return the time of the last call to Read()
func (self *TimedReadCloser) LastReadAt() time.Time {
	return self.lastRead
}

// Return the duration of time since the last Read() occurred.
func (self *TimedReadCloser) SinceLastRead() time.Duration {
	if !self.lastRead.IsZero() {
		return time.Since(self.lastRead)
	} else {
		return 0
	}
}

// Return a running duration of how long reading has been happening.  Updated on
// every call to Read().
func (self *TimedReadCloser) ReadDuration() time.Duration {
	return self.readTook
}

func (self *TimedReadCloser) CloseDuration() time.Duration {
	return self.closeTook
}

func (self *TimedReadCloser) Duration() time.Duration {
	return self.took
}

func (self *TimedReadCloser) BytesRead() int64 {
	return self.bytesRead
}
