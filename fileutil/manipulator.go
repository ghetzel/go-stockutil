package fileutil

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// By default, the underlying bufio.Scanner that tokenizes the input data
// will discard the tokens that it's splitting on.  however, in most cases, we
// don't actually want this.  however, if the thing we're removing IS this token,
// then we will have removed the token, then immediately put it back in the stream.
//
// Returning the ErrSkipToken error will tell the ReadManipulator to not put this token
// back into the stream, but otherwise not produce an actual error during read.
var ErrSkipToken = errors.New(`skip token`)

type ReadManipulatorFunc func(data []byte) ([]byte, error)

type ReadManipulator struct {
	reader    io.Reader
	fn        ReadManipulatorFunc
	splitter  bufio.SplitFunc
	scanner   *bufio.Scanner
	buffer    *bytes.Buffer
	lastToken []byte
}

func NewReadManipulator(reader io.Reader, fns ...ReadManipulatorFunc) *ReadManipulator {
	var rm = &ReadManipulator{
		reader:   reader,
		splitter: bufio.ScanLines,
		buffer:   bytes.NewBuffer(nil),
	}

	if len(fns) > 0 && fns[0] != nil {
		rm.fn = fns[0]
	}

	return rm
}

func (manip *ReadManipulator) Split(split bufio.SplitFunc) {
	manip.splitter = split
}

func (manip *ReadManipulator) Read(b []byte) (int, error) {
	if manip.fn != nil {
		// initialize the scanner if we need to
		if manip.scanner == nil {
			manip.scanner = bufio.NewScanner(manip.reader)
			manip.scanner.Split(manip.interceptToken)
		}

		// if there's more to scan...
		for manip.scanner.Scan() {
			var data = manip.scanner.Bytes()

			// get the scanned bytes, run them through the manip. function
			if out, err := manip.fn(data); err == nil || err == ErrSkipToken {
				if err == nil {
					out = append(out, manip.lastToken...)
				}

				manip.lastToken = nil

				// write the manipulated bytes to the buffer
				if n, err := manip.buffer.Write(out); err != nil {
					return n, err
				}

				// loop until we've put enough data in the buffer to satisfy the
				// requested read
				if manip.buffer.Len() >= len(b) {
					break
				}
			} else {
				return 0, err
			}
		}

		// check for scan errors
		if err := manip.scanner.Err(); err != nil {
			return 0, err
		}

		// return whats in the buffer, and keep doing this until its empty
		return manip.buffer.Read(b)
	} else {
		return manip.reader.Read(b)
	}
}

func (manip *ReadManipulator) Close() error {
	if manip.scanner != nil {
		manip.scanner = nil
	}

	manip.lastToken = nil
	manip.buffer = bytes.NewBuffer(nil)

	if closer, ok := manip.reader.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}

func (manip *ReadManipulator) interceptToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = manip.splitter(data, atEOF)

	if advance > 0 && len(data) >= advance {
		manip.lastToken = append(manip.lastToken, data[len(token):advance]...)
	}

	return
}
