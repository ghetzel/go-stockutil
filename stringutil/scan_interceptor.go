package stringutil

import (
	"bufio"
	"bytes"
)

type InterceptFunc func(seq []byte)

// A ScanInterceptor is used as a SplitFunc on a bufio.Scanner.  It will look at the stream of bytes being scanned for
// specific substrings.  The registered handler function associated with a substring will be called whenever it is seen
// in the stream.  The passthrough SplitFunc is called as normal.  This allows for a stream to be
// split and processed while also being inspected for specific content, allowing the user to react to that content
// as it comes by.
type ScanInterceptor struct {
	Disabled           bool
	accumulator        *bytes.Buffer
	subsequences       map[string]InterceptFunc
	interceptStats     map[string]int64
	longestSubsequence int
	totalWritten       int64
	totalAdvanced      int64
	highWaterMark      map[string]int64
	passthrough        bufio.SplitFunc
}

func NewScanInterceptor(passthrough bufio.SplitFunc, intercepts ...map[string]InterceptFunc) *ScanInterceptor {
	var intercept map[string]InterceptFunc

	if len(intercepts) == 0 {
		intercept = make(map[string]InterceptFunc)
	} else {
		intercept = intercepts[0]
	}

	// return a new, empty interceptor
	return &ScanInterceptor{
		passthrough:    passthrough,
		accumulator:    bytes.NewBuffer(nil),
		subsequences:   intercept,
		highWaterMark:  make(map[string]int64),
		interceptStats: make(map[string]int64),
	}
}

// Add an intercept sequence and handler.  If the sequence is already registered, its handler
// function will be replaced with this one.
func (self *ScanInterceptor) Intercept(sequence string, handler InterceptFunc) {
	self.subsequences[sequence] = handler

	for k, _ := range self.subsequences {
		if len(k) > self.longestSubsequence {
			self.longestSubsequence = len(k)
		}
	}
}

// Implements the bufio.SplitFunc function signature for use in a bufio.Scanner.
func (self *ScanInterceptor) Scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// call the SplitFunc we were given
	advance, token, err = self.passthrough(data, atEOF)

	if !self.Disabled {
		if len(token) > 0 {
			if _, err := self.accumulator.Write(token); err != nil {
				return 0, nil, err
			}

			// however far we just advanced (if at all), keep track of that
			self.totalWritten += int64(len(token))
		}

		// if we've accumulated *at least* as many bytes as our longest subsequence, then
		// we go to work...
		if processedLen := self.accumulator.Len(); processedLen >= self.longestSubsequence {
			// get the bytes we've accumulated since start or the last time we reset
			var soFar = self.accumulator.Bytes()

			// for each registered subsequence...
			for k, handler := range self.subsequences {
				var subseq = []byte(k)

				// skip zero-length matches
				if len(subseq) == 0 {
					continue
				}

				// the High Water Mark (HWM) represents the furthest we've ever gotten in the stream.
				// we make sure that our current HWM is *before* the end of the stream, so that if this
				// SplitFunc is called repeatedy for the same data (which can happen), we're not firing off
				// multiple handler calls for the same position(s).
				//
				if self.highWaterMark[k] > self.totalWritten {
					continue
				}

				var work = soFar

				// find the index in the stream of our match (if any)
				for {
					if indexOf := bytes.Index(work, subseq); indexOf >= 0 {
						// mark the end of the stream (so we ensure we dont fire events for anything before this point)
						var endIndex = indexOf + len(subseq)

						// fire the handler
						handler(work[indexOf:endIndex])
						work = work[endIndex:]

						self.interceptStats[k] = self.interceptStats[k] + 1

						// advance the HWM for this interceptor past this result
						self.highWaterMark[k] = self.totalWritten + int64(endIndex) + 1
					} else {
						break
					}
				}
			}

			// reset the accumulator, we go again!
			self.accumulator = bytes.NewBuffer(nil)
		}
	}

	// return the results of the SplitFunc we were given
	return advance, token, err
}

// Return the total number of bytes this scanner has scanned.
func (self *ScanInterceptor) BytesScanned() int64 {
	return self.totalWritten
}

// Returns a map of intercept sequences and the number of times each one was fired.
func (self *ScanInterceptor) InterceptCounts() map[string]int64 {
	return self.interceptStats
}
