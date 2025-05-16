package fileutil

import (
	"io"

	"github.com/ghetzel/go-stockutil/log"
)

type CloserFunc = func(io.ReadCloser) error

type PostReadCloser struct {
	upstream io.ReadCloser
	closer   func(io.ReadCloser) error
}

// Implements an io.ReadCloser that can be configured to perform cleanup options whenever the
// Close() function is called.  If CloserFunc is non-nil, it will be given the upstream ReadCloser
// as an argument and will be responsible for calling Close() on it.  If nil, upstream's Close()
// function will be called directly on Close.
func NewPostReadCloser(upstream io.ReadCloser, closer CloserFunc) *PostReadCloser {
	return &PostReadCloser{
		upstream: upstream,
		closer:   closer,
	}
}

func (postcloser *PostReadCloser) Read(b []byte) (int, error) {
	return postcloser.upstream.Read(b)
}

func (postcloser *PostReadCloser) Close() error {
	if fn := postcloser.closer; fn != nil {
		return fn(postcloser.upstream)
	} else {
		return postcloser.upstream.Close()
	}
}

type MultiCloser struct {
	closers []io.Closer
}

func NewMultiCloser(closers ...io.Closer) *MultiCloser {
	return &MultiCloser{
		closers: closers,
	}
}

func (multicloser *MultiCloser) Close() (merr error) {
	for _, closer := range multicloser.closers {
		if closer != nil {
			merr = log.AppendError(merr, closer.Close())
		}
	}

	return
}
