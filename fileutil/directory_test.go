package fileutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ghetzel/testify/require"
)

type appendCoolWriter struct {
	writer io.Writer
}

func (appender *appendCoolWriter) Write(b []byte) (int, error) {
	_, err := appender.writer.Write([]byte(string(b) + "cool.\n"))
	return len(b), err
}

func TestDirReader(t *testing.T) {
	var assert = require.New(t)

	var dread = NewDirReader(`./testdir`)
	defer dread.Close()
	data, err := io.ReadAll(dread)
	assert.NoError(err)
	assert.Len(data, 22)
	assert.Equal("a\nb\nc\nd1\nd2\nd3\ne11\ne2\n", string(data))

	// close and see if DirReader reset properly
	assert.NoError(dread.Close())

	data, err = io.ReadAll(dread)
	assert.NoError(err)
	assert.Len(data, 22)
	assert.Equal("a\nb\nc\nd1\nd2\nd3\ne11\ne2\n", string(data))

	dread = NewDirReader(`./testdir/d`)
	data, err = io.ReadAll(dread)
	assert.NoError(err)
	assert.Len(data, 9)
	assert.Equal("d1\nd2\nd3\n", string(data))
}

func TestDirReaderSkipFunc(t *testing.T) {
	var assert = require.New(t)

	var dread = NewDirReader(`./testdir`)
	defer dread.Close()
	dread.SetSkipFunc(func(p string) bool {
		var filename = strings.TrimSuffix(p, filepath.Ext(p))

		t.Logf("%s: %v", filename, strings.HasSuffix(filename, `1`))

		return strings.HasSuffix(filename, `1`)
	})

	data, err := io.ReadAll(dread)
	assert.NoError(err)
	assert.Len(data, 15)
	assert.Equal("a\nb\nc\nd2\nd3\ne2\n", string(data))
}

func TestCopyDir(t *testing.T) {
	var assert = require.New(t)

	var buf = bytes.NewBuffer(nil)
	var cool = &appendCoolWriter{
		writer: buf,
	}

	assert.NoError(CopyDir(`./testdir`, func(path string, info os.FileInfo, err error) (io.Writer, error) {
		if info.IsDir() {
			return nil, nil
		} else {
			return cool, err
		}
	}))

	assert.Equal(buf.Len(), 70)
	assert.EqualValues(
		"a\ncool.\nb\ncool.\nc\ncool.\nd1\ncool.\nd2\ncool.\nd3\ncool.\ne11\ncool.\ne2\ncool.\n",
		buf.String(),
	)
}
