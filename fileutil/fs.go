// Helpers for working with files and filesystems
package fileutil

import (
	"net/http"
	"os"
	"regexp"
)

type RewriteFileSystem struct {
	FileSystem http.FileSystem
	Find       *regexp.Regexp
	Replace    string
	MustMatch  bool
}

func (rwfs RewriteFileSystem) Open(name string) (http.File, error) {
	if rwfs.Find != nil {
		if rwfs.Find.MatchString(name) {
			name = rwfs.Find.ReplaceAllString(name, rwfs.Replace)
		} else if rwfs.MustMatch {
			return nil, os.ErrNotExist
		}
	}

	return rwfs.FileSystem.Open(name)
}
