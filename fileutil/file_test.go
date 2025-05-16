package fileutil

import (
	"io"
	"os"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestSetExt(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(``, SetExt(``, ``))
	assert.Equal(`/nothingburger.txt`, SetExt(`/nothingburger.txt`, ``))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `.jpg`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `jpg`))
	assert.Equal(`/nothingburger.txt`, SetExt(`/nothingburger.txt`, `.jpg`, `.bmp`))
	assert.Equal(`/nothingburger.txt`, SetExt(`/nothingburger.txt`, `jpg`, `.bmp`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `.jpg`, `.txt`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `jpg`, `.txt`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `.jpg`, `txt`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `jpg`, `txt`))

	assert.Equal(`/nothingburger.info.xml`, SetExt(`/nothingburger.info.json`, `xml`))
	assert.Equal(`/nothingburger.info.xml`, SetExt(`/nothingburger.info.json`, `.xml`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `.xml`, `info.json`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `.xml`, `.info.json`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `xml`, `info.json`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `xml`, `.info.json`))
}

func reset(s io.Seeker) error {
	_, err := s.Seek(0, 0)
	return err
}

func TestSameFile(t *testing.T) {
	var assert = require.New(t)

	assert.True(SameFile(`file_test.go`, `file_test.go`))
	assert.True(SameFile(`file.go`, `file.go`))

	i, err := os.Stat(`file_test.go`)
	assert.NoError(err)
	o, err := os.Stat(`file.go`)
	assert.NoError(err)

	assert.True(SameFile(i, i))
	assert.True(SameFile(o, o))
	assert.False(SameFile(i, o))
	assert.False(SameFile(o, i))
	assert.False(SameFile(i, nil))
	assert.False(SameFile(nil, i))
	assert.False(SameFile(o, nil))
	assert.False(SameFile(nil, o))

	assert.True(SameFile(i, `file_test.go`))
	assert.True(SameFile(`file_test.go`, i))
	assert.True(SameFile(o, `file.go`))
	assert.True(SameFile(`file.go`, o))

	f1, err := os.Open(`file_test.go`)
	assert.NoError(err)
	defer f1.Close()

	f2, err := os.Open(`file.go`)
	assert.NoError(err)
	defer f2.Close()

	assert.True(SameFile(f1, f1))
	assert.NoError(reset(f1))

	assert.True(SameFile(f1, f1))
	assert.NoError(reset(f1))

	assert.True(SameFile(i, f1))
	assert.NoError(reset(f1))

	assert.False(SameFile(f1, f2))
	assert.NoError(reset(f1))
	assert.NoError(reset(f2))

	assert.False(SameFile(f2, f1))
	assert.NoError(reset(f1))
	assert.NoError(reset(f2))
}

func TestIsHiddenFile(t *testing.T) {
	var assert = require.New(t)

	assert.True(IsHiddenFile(`.hidden`))
	assert.True(IsHiddenFile(`path/to/.hidden`))
	assert.False(IsHiddenFile(``))
	assert.False(IsHiddenFile(`.`))
	assert.False(IsHiddenFile(`..`))
	assert.False(IsHiddenFile(`visible`))
	assert.False(IsHiddenFile(`path/to/visible`))
}
