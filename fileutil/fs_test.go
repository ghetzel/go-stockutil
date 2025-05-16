package fileutil

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/ghetzel/testify/require"
)

type TestFileSystem map[string]http.File

func (testfs TestFileSystem) Open(name string) (http.File, error) {
	fmt.Printf("Opening %q\n", name)

	if file, ok := testfs[name]; ok {
		return file, nil
	}

	return nil, os.ErrNotExist
}

func TestRewriteFileSystem(t *testing.T) {
	var assert = require.New(t)

	// ------------------------------------------------------------------------
	var rwfs = RewriteFileSystem{
		FileSystem: TestFileSystem{
			`/test`: nil,
		},
	}

	_, err := rwfs.Open(`/test`)
	assert.Nil(err)

	// ------------------------------------------------------------------------
	rwfs = RewriteFileSystem{
		FileSystem: TestFileSystem{
			`/test`: nil,
		},
		Find: regexp.MustCompile(`^/strip`),
	}

	_, err = rwfs.Open(`/test`)
	assert.Nil(err)

	_, err = rwfs.Open(`/strip/test`)
	assert.Nil(err)

	// ------------------------------------------------------------------------
	rwfs = RewriteFileSystem{
		FileSystem: TestFileSystem{
			`/test`: nil,
		},
		Find:      regexp.MustCompile(`^/strip`),
		MustMatch: true,
	}

	_, err = rwfs.Open(`/test`)
	assert.Equal(os.ErrNotExist, err)

	_, err = rwfs.Open(`/strip/test`)
	assert.Nil(err)

	// ------------------------------------------------------------------------
	rwfs = RewriteFileSystem{
		FileSystem: TestFileSystem{
			`/other/test`: nil,
		},
		Find:      regexp.MustCompile(`^/strip`),
		Replace:   `/other`,
		MustMatch: true,
	}

	_, err = rwfs.Open(`/strip/test`)
	assert.Nil(err)

	_, err = rwfs.Open(`/other/test`)
	assert.Equal(os.ErrNotExist, err)

	// ------------------------------------------------------------------------
	rwfs = RewriteFileSystem{
		FileSystem: TestFileSystem{
			`/before/after/test`: nil,
		},
		Find:    regexp.MustCompile(`^/(?P<first>[^/]+)/(?P<second>[^/]+)`),
		Replace: `/${second}/${first}`,
	}

	_, err = rwfs.Open(`/after/before/test`)
	assert.Nil(err)

	_, err = rwfs.Open(`/before/after/test`)
	assert.Equal(os.ErrNotExist, err)
}
