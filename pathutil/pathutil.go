// Additional helpers for working with file paths and filesystem information
package pathutil

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/gobwas/glob"
)

// Returns whether a given path matches a glob pattern.
//
// via github.com/gobwas/glob:
//
// Compile creates Glob for given pattern and strings (if any present after pattern) as separators.
// The pattern syntax is:
//
//	pattern:
//	    { term }
//
//	term:
//	    `*`         matches any sequence of non-separator characters
//	    `**`        matches any sequence of characters
//	    `?`         matches any single non-separator character
//	    `[` [ `!` ] { character-range } `]`
//	                character class (must be non-empty)
//	    `{` pattern-list `}`
//	                pattern alternatives
//	    c           matches character c (c != `*`, `**`, `?`, `\`, `[`, `{`, `}`)
//	    `\` c       matches character c
//
//	character-range:
//	    c           matches character c (c != `\\`, `-`, `]`)
//	    `\` c       matches character c
//	    lo `-` hi   matches character c for lo <= c <= hi
//
//	pattern-list:
//	    pattern { `,` pattern }
//	                comma-separated (without spaces) patterns
func MatchPath(pattern string, path string) bool {
	if g, err := glob.Compile(pattern); err == nil {
		return g.Match(path)
	}

	return false
}

// ExpandUser replaces the tilde (~) in a path into the current user's home directory.
func ExpandUser(path string) (string, error) {
	if u, err := user.Current(); err == nil {
		var fullTilde = fmt.Sprintf("~%s", u.Name)

		if strings.HasPrefix(path, `~/`) || path == `~` {
			return strings.Replace(path, `~`, u.HomeDir, 1), nil
		}

		if strings.HasPrefix(path, fullTilde+`/`) || path == fullTilde {
			return strings.Replace(path, fullTilde, u.HomeDir, 1), nil
		}

		return path, nil
	} else {
		return path, err
	}
}

// Returns true if the given path is a regular file, is executable by any user, and has a non-zero size.
func IsNonemptyExecutableFile(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.Size() > 0 && (stat.Mode().Perm()&0111) != 0 {
		return true
	}

	return false
}

// Returns true if the given path is a regular file with a non-zero size.
func IsNonemptyFile(path string) bool {
	if FileExists(path) {
		if stat, err := os.Stat(path); err == nil && stat.Size() > 0 {
			return true
		}
	}

	return false
}

// Returns true if the given path is a directory with items in it.
func IsNonemptyDir(path string) bool {
	if DirExists(path) {
		if entries, err := os.ReadDir(path); err == nil && len(entries) > 0 {
			return true
		}
	}

	return false
}

// Returns true if the given path exists.
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

// Returns true if the given path exists and is a symbolic link.
func LinkExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return IsSymlink(stat.Mode())
	}

	return false
}

// Returns true if the given path exists and is a regular file.
func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.Mode().IsRegular()
	}

	return false
}

// Returns true if the given path exists and is a directory.
func DirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}

	return false
}

func IsSymlink(mode os.FileMode) bool {
	return (mode&os.ModeSymlink != 0)
}

func IsDevice(mode os.FileMode) bool {
	return (mode&os.ModeDevice != 0)
}

func IsCharDevice(mode os.FileMode) bool {
	return (mode&os.ModeCharDevice != 0)
}

func IsNamedPipe(mode os.FileMode) bool {
	return (mode&os.ModeNamedPipe != 0)
}

func IsSocket(mode os.FileMode) bool {
	return (mode&os.ModeSocket != 0)
}

func IsSticky(mode os.FileMode) bool {
	return (mode&os.ModeSticky != 0)
}

func IsSetuid(mode os.FileMode) bool {
	return (mode&os.ModeSetuid != 0)
}

func IsSetgid(mode os.FileMode) bool {
	return (mode&os.ModeSetgid != 0)
}

func IsTemporary(mode os.FileMode) bool {
	return (mode&os.ModeTemporary != 0)
}

func IsExclusive(mode os.FileMode) bool {
	return (mode&os.ModeExclusive != 0)
}

func IsAppend(mode os.FileMode) bool {
	return (mode&os.ModeAppend != 0)
}

// Returns true if the given file can be opened for reading by the current user.
func IsReadable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_RDONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

// Returns true if the given file can be opened for writing by the current user.
func IsWritable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

// Returns true if the given file can be opened for appending by the current user.
func IsAppendable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_APPEND, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}
