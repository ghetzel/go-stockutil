package fileutil

import (
	"io"
	"os"
	"path/filepath"
	"sort"
)

type DirReaderOption int

const (
	NoRecursive DirReaderOption = iota
	FailOnError
	IncludeHidden
)

type DirReaderOptions []DirReaderOption

func (dir DirReaderOptions) Has(option DirReaderOption) bool {
	for _, opt := range dir {
		if opt == option {
			return true
		}
	}

	return false
}

type SkipFunc func(string) bool

// A DirReader provides a streaming io.Reader interface to all files in a given
// directory, with options for handling unreadable entries and recursion.
type DirReader struct {
	root         string
	options      DirReaderOptions
	loaded       bool
	entries      []os.FileInfo
	size         int64
	currentEntry int
	current      io.ReadCloser
	skipFn       SkipFunc
}

func NewDirReader(path string, options ...DirReaderOption) *DirReader {
	return &DirReader{
		root:    path,
		options: DirReaderOptions(options),
	}
}

// Set a function that will be called for each path encountered while reading.
// If this function returns true, that path (and its descedants) will not be read.
func (dir *DirReader) SetSkipFunc(fn SkipFunc) {
	dir.skipFn = fn
}

func (dir *DirReader) setup() error {
	if rt, entries, err := dir.readDir(); err == nil {
		dir.root = rt
		dir.current = nil
		dir.currentEntry = 0
		dir.entries = entries
		dir.size = 0
		dir.loaded = true
		return nil
	} else {
		return err
	}
}

func (dir *DirReader) advanceAndRead(b []byte) (int, error) {
	if dir.current != nil {
		if err := dir.current.Close(); err != nil {
			return 0, err
		}
	}

	// if the current entry is in-bounds
	if dir.currentEntry < len(dir.entries) {
		var entry = dir.entries[dir.currentEntry]
		dir.currentEntry += 1
		var path = filepath.Join(dir.root, entry.Name())

		// if a skipFunc we provided, and it returned false, return a skiperr from here
		if skipFn := dir.skipFn; skipFn != nil && skipFn(path) {
			return dir.advanceAndRead(b)
		}

		if !dir.options.Has(IncludeHidden) && IsHiddenFile(path) {
			return dir.advanceAndRead(b)
		}

		if entry.IsDir() {
			if !dir.options.Has(NoRecursive) {
				// directories that we recurse into are just instances of DirReaders whose root is the directory
				var subreader = NewDirReader(path, dir.options...)
				subreader.SetSkipFunc(dir.skipFn)

				dir.current = subreader
				return dir.current.Read(b)
			} else {
				return dir.advanceAndRead(b)
			}
		} else if file, err := os.Open(path); err == nil {
			dir.current = file
			return dir.current.Read(b)
		} else {
			return 0, err
		}
	} else {
		return 0, io.EOF
	}
}

func (dir *DirReader) Read(b []byte) (int, error) {
	// do initial setup if we're starting out
	if !dir.loaded {
		if err := dir.setup(); err != nil {
			return 0, err
		}
	}

	// check if we have a current file
	if dir.current != nil {
		if n, err := dir.current.Read(b); err == nil {
			// if so, read from that file
			return n, nil
		} else if err == io.EOF || !dir.options.Has(FailOnError) {
			// if the current file is EOF (or we're skipping errors), advance to the next one and start reading it
			// keep advancing until the error is not skipEntryErr
			return dir.advanceAndRead(b)
		} else {
			return n, err
		}
	} else {
		return dir.advanceAndRead(b)
	}
}

// close open files and reset the internal reader
func (dir *DirReader) Close() error {
	if dir.current != nil {
		dir.current.Close()
	}

	dir.loaded = false

	return nil
}

// read the immediate entries from  the current root directory, or if the current root
// is a file, treat it like a directory of one entry
func (dir *DirReader) readDir() (string, []os.FileInfo, error) {
	if root, err := ExpandUser(dir.root); err == nil {
		if DirExists(root) {
			var dentries, err = os.ReadDir(dir.root)
			var entries []os.FileInfo

			for _, de := range dentries {
				if fi, err := de.Info(); err == nil {
					entries = append(entries, fi)
				} else {
					return ``, nil, err
				}
			}

			sort.Slice(entries, func(i int, j int) bool {
				return entries[i].Name() < entries[j].Name()
			})

			return root, entries, err
		} else if stat, err := os.Stat(root); err == nil {
			return root, []os.FileInfo{stat}, nil
		} else {
			return ``, nil, err
		}
	} else {
		return ``, nil, err
	}
}

type CopyEntryFunc func(path string, info os.FileInfo, err error) (io.Writer, error)

// Recursively walk the entries of a given directory, calling CopyEntryFunc for each entry.
// The io.Writer returned from the function will have that file's contents written to it.  If
// the io.Writer is nil, the file will not be written anywhere but no error will be returned.
// If CopyEntryFunc returns an error, the behavior will be consistent with filepath.WalkFunc
func CopyDir(root string, fn CopyEntryFunc) error {
	return copyDir(root, false, fn)
}

// A version of CopyDir that includes hidden files
func CopyDirWithHidden(root string, fn CopyEntryFunc) error {
	return copyDir(root, true, fn)
}

func copyDir(root string, withHidden bool, fn CopyEntryFunc) error {
	if p, err := ExpandUser(root); err == nil {
		root = p
	} else {
		return err
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !withHidden && IsHiddenFile(path) {
			return nil
		}

		if w, err := fn(path, info, err); err == nil && w != nil {
			if file, err := os.Open(path); err == nil {
				defer file.Close()
				_, err := io.Copy(w, file)
				return err
			} else {
				return err
			}
		} else if err != nil {
			return err
		} else {
			return nil
		}
	})
}
