package fileutil

import (
	"container/list"
	"io"
	"sync"

	"github.com/ghetzel/go-stockutil/log"
)

type ExtendableReader struct {
	sources list.List
	srclock sync.Mutex
	current io.ReadCloser
}

func (xreader *ExtendableReader) AppendSource(rc io.ReadCloser) {
	if rc != nil {
		xreader.srclock.Lock()
		xreader.sources.PushBack(rc)
		xreader.srclock.Unlock()
	}
}

func (xreader *ExtendableReader) Read(b []byte) (int, error) {
	if err := xreader.closeAndAdvanceSources(); err != nil {
		return 0, err
	}

	if xreader.current == nil {
		return 0, io.EOF
	} else if n, err := xreader.current.Read(b); err == nil {
		return n, nil
	} else if err == io.EOF {
		return xreader.Read(b)
	} else {
		return n, err
	}
}

func (xreader *ExtendableReader) Close() error {
	var merr error

	xreader.srclock.Lock()
	defer xreader.srclock.Unlock()

	for e := xreader.sources.Front(); e != nil; e = e.Next() {
		merr = log.AppendError(merr, e.Value.(io.ReadCloser).Close())
		xreader.sources.Remove(e)
	}

	return merr
}

func (xreader *ExtendableReader) closeAndAdvanceSources() error {
	if xreader.current != nil {
		if err := xreader.current.Close(); err != nil {
			return err
		}
	}

	xreader.srclock.Lock()
	defer xreader.srclock.Unlock()

	if xreader.sources.Len() > 0 {
		if el := xreader.sources.Front(); el != nil {
			xreader.current = el.Value.(io.ReadCloser)
			xreader.sources.Remove(el)

			return nil
		}
	}

	return io.EOF
}
