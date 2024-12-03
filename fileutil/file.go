package fileutil

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghetzel/go-stockutil/convutil"
	"github.com/ghetzel/go-stockutil/pathutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	isatty "github.com/mattn/go-isatty"
)

// Alias functions from pathutil as a convenience
var DirExists = pathutil.DirExists
var Exists = pathutil.Exists
var ExpandUser = pathutil.ExpandUser
var FileExists = pathutil.FileExists
var IsAppend = pathutil.IsAppend
var IsAppendable = pathutil.IsAppendable
var IsCharDevice = pathutil.IsCharDevice
var IsDevice = pathutil.IsDevice
var IsExclusive = pathutil.IsExclusive
var IsNamedPipe = pathutil.IsNamedPipe
var IsNonemptyDir = pathutil.IsNonemptyDir
var IsNonemptyExecutableFile = pathutil.IsNonemptyExecutableFile
var IsNonemptyFile = pathutil.IsNonemptyFile
var IsReadable = pathutil.IsReadable
var IsSetgid = pathutil.IsSetgid
var IsSetuid = pathutil.IsSetuid
var IsSocket = pathutil.IsSocket
var IsSticky = pathutil.IsSticky
var IsSymlink = pathutil.IsSymlink
var IsTemporary = pathutil.IsTemporary
var IsWritable = pathutil.IsWritable
var LinkExists = pathutil.LinkExists
var CompareHasher hash.Hash = sha256.New()

func MustExpandUser(path string) string {
	if expanded, err := ExpandUser(path); err == nil {
		return expanded
	} else {
		panic(err.Error())
	}
}

// Return true if the given FileInfo sports a ModTime later than the current file.
func IsModifiedAfter(stat os.FileInfo, current string) bool {
	if Exists(current) {
		current = MustExpandUser(current)

		if currentStat, err := os.Stat(current); err == nil {
			// return whether the new file is modified AFTER the current one
			return stat.ModTime().After(currentStat.ModTime())
		} else {
			// fail closed; if we can't stat the existing file then don't assert that we know better
			return false
		}
	} else {
		// if the current file does not exist, then whatever file we have is newer
		return true
	}
}

func IsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// Takes the given string filename, []byte, io.Reader, or io.ReadCloser and returns
// the bytes therein.
func ReadAll(file interface{}) ([]byte, error) {
	var rc io.ReadCloser

	if b, ok := file.([]byte); ok {
		return b, nil
	} else if r, ok := file.(io.ReadCloser); ok {
		rc = r
	} else if r, ok := file.(io.Reader); ok {
		rc = ioutil.NopCloser(r)
	} else if filename, ok := file.(string); ok {
		if f, err := os.Open(filename); err == nil {
			rc = f
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("must provide one of: string, []byte, io.Reader, or io.ReadCloser")
	}

	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// A string version of ReadAll.
func ReadAllString(file interface{}) (string, error) {
	if data, err := ReadAll(file); err == nil {
		return string(data), nil
	} else {
		return ``, err
	}
}

// Attempt to call ReadAllString, but will return an empty string if there is an error.  Does not panic.
func Cat(file interface{}) string {
	if data, err := ReadAllString(file); err == nil {
		return data
	} else {
		return ``
	}
}

// Read all lines of text from the given file and return them as a slice.
func ReadAllLines(file interface{}) ([]string, error) {
	if data, err := ReadAllString(file); err == nil {
		return strings.Split(data, "\n"), nil
	} else {
		return nil, err
	}
}

// A panicky version of ReadAllLines.
func MustReadAllLines(file interface{}) []string {
	if lines, err := ReadAllLines(file); err == nil {
		return lines
	} else {
		panic(err.Error())
	}
}

// Attempts to call ReadAllLines, but will return an empty slice if there is an error.  Does not panic.
func ShouldReadAllLines(file interface{}) []string {
	if lines, err := ReadAllLines(file); err == nil {
		return lines
	} else {
		return nil
	}
}

// Attempt to return the nth line (starting from 1) in the given file or reader.
func GetNthLine(file interface{}, number int) (string, error) {
	if lines, err := ReadAllLines(file); err == nil {
		if number < 1 {
			return ``, fmt.Errorf("line number must be >= 1")
		}

		if number <= len(lines) {
			return lines[number-1], nil
		} else {
			return ``, fmt.Errorf("line exceeds file line count")
		}
	} else {
		return ``, err
	}
}

// A panicky version of GetNthLine.
func MustGetNthLine(file interface{}, number int) string {
	if line, err := GetNthLine(file, number); err == nil {
		return line
	} else {
		panic(err.Error())
	}
}

// Attempts to call GetNthLine, but will return an empty string if there is an error.  Does not panic.
func ShouldGetNthLine(file interface{}, number int) string {
	if line, err := GetNthLine(file, number); err == nil {
		return line
	} else {
		return ``
	}
}

func ReadFirstLine(file interface{}) (string, error) {
	if lines, err := ReadAllLines(file); err == nil {
		return lines[0], nil
	} else {
		return ``, err
	}
}

func MustReadAll(file interface{}) []byte {
	if data, err := ReadAll(file); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}

func MustReadAllString(file interface{}) string {
	if data, err := ReadAllString(file); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}

func ChecksumFile(filename string, fn interface{}) ([]byte, error) {
	var hash crypto.Hash

	if h, ok := fn.(crypto.Hash); ok {
		hash = h
	} else {
		switch strings.ToLower(typeutil.String(fn)) {
		case `md4`:
			hash = crypto.MD4
		case `md5`:
			hash = crypto.MD5
		case `sha1`:
			hash = crypto.SHA1
		case `sha224`:
			hash = crypto.SHA224
		case `sha256`:
			hash = crypto.SHA256
		case `sha384`:
			hash = crypto.SHA384
		case `sha512`:
			hash = crypto.SHA512
		case `md5sha1`:
			hash = crypto.MD5SHA1
		case `ripemd160`:
			hash = crypto.RIPEMD160
		case `sha3_224`:
			hash = crypto.SHA3_224
		case `sha3_256`:
			hash = crypto.SHA3_256
		case `sha3_384`:
			hash = crypto.SHA3_384
		case `sha3_512`:
			hash = crypto.SHA3_512
		case `sha512_224`:
			hash = crypto.SHA512_224
		case `sha512_256`:
			hash = crypto.SHA512_256
		case `blake2s_256`:
			hash = crypto.BLAKE2s_256
		case `blake2b_256`:
			hash = crypto.BLAKE2b_256
		case `blake2b_384`:
			hash = crypto.BLAKE2b_384
		case `blake2b_512`:
			hash = crypto.BLAKE2b_512
		default:
			return nil, fmt.Errorf("Unknown hash function %q", fn)
		}
	}

	if hash.Available() {
		hasher := hash.New()

		if file, err := os.Open(filename); err == nil {
			defer file.Close()

			if _, err := io.Copy(hasher, file); err == nil {
				sum := hasher.Sum(nil)
				return sum, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Hash function %v is not available", hash)
	}
}

// Write the contents of the given io.Reader, []byte, or string to the specified filename.
// Filename paths containing tilde (~) will automatically expand to the current
// user's home directory, and all intermediate parent directories will be automatically
// created.  Will return the number of bytes written, or an error.
func WriteFile(input interface{}, filename string) (int64, error) {
	var reader io.Reader

	if r, ok := input.(io.Reader); ok {
		reader = r
	} else if b, ok := input.([]byte); ok {
		reader = bytes.NewBuffer(b)
	} else if s, ok := input.(string); ok {
		reader = bytes.NewBufferString(s)
	} else {
		return 0, fmt.Errorf("Cannot use %T as input", input)
	}

	if expanded, err := ExpandUser(filename); err == nil {
		parent := filepath.Dir(expanded)

		// create parent directory automatically
		if !DirExists(parent) {
			if err := os.MkdirAll(parent, 0700); err != nil {
				return 0, err
			}
		}

		// open the destination file for writing
		if dest, err := os.Create(expanded); err == nil {
			defer dest.Close()
			return io.Copy(dest, reader)
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

// Same as WriteFile, but will panic if the file cannot be written.
func MustWriteFile(input interface{}, filename string) int64 {
	if n, err := WriteFile(input, filename); err == nil {
		return n
	} else {
		panic(err.Error())
	}
}

// Same as WriteFile, but writes the given input to a temporary file, returning
// the filename.
func WriteTempFile(input interface{}, pattern string) (string, error) {
	if tmp, err := ioutil.TempFile(``, pattern); err == nil {
		defer tmp.Close()

		if _, err := WriteFile(input, tmp.Name()); err == nil {
			return tmp.Name(), nil
		} else {
			defer os.Remove(tmp.Name())
			return ``, err
		}
	} else {

		return ``, err
	}
}

// Same as MustWriteFile, but writes the given input to a temporary file, returning
// the filename.  The function will panic if the file cannot be written.
func MustWriteTempFile(input interface{}, prefix string) string {
	if n, err := WriteTempFile(input, prefix); err == nil {
		return n
	} else {
		panic(err.Error())
	}
}

// Write the contents of the given reader to the specified filename.
// Filename paths containing tilde (~) will automatically expand to the current
// user's home directory, and all intermediate parent directories will be automatically
// created.  Will return the number of bytes written, or an error.

// Copy a file from one place to another.  Source can be an io.Reader or string.  If source is a
// string, the string will be passed to the Open() function as a URL.  Destination can be an
// io.Writer or string.  If destination is a string, it will be treated as a local filesystem path
// to write the data read from source to.
//
// If either source or destination implements io.Closer, thee files will be closed before this
// function returns.
func CopyFile(source interface{}, destination interface{}) error {
	var sreader io.Reader
	var dwriter io.Writer

	// open or otherwise get the source
	if sfilename, ok := source.(string); ok {
		if sr, err := Open(sfilename); err == nil {
			sreader = sr
		} else {
			return err
		}
	} else if sr, ok := source.(io.Reader); ok {
		sreader = sr
	} else {
		return fmt.Errorf("Unsupported source %T", source)
	}

	// open the destination for writing
	if dfilename, ok := source.(string); ok {
		if dfile, err := os.Create(dfilename); err == nil {
			dwriter = dfile
		} else {
			return err
		}
	} else if dw, ok := destination.(io.Writer); ok {
		dwriter = dw
	} else {
		return fmt.Errorf("Unsupported source %T", source)
	}

	// defer closing of both source and destination if they support it
	defer func() {
		if sc, ok := sreader.(io.Closer); ok {
			sc.Close()
		}

		if dc, ok := dwriter.(io.Closer); ok {
			dc.Close()
		}
	}()

	// copy and return
	_, err := io.Copy(dwriter, sreader)
	return err
}

// Detects the file extension of the given path and replaces it with the given extension.  The optional
// second argument allows you to explictly specify the extension (if known).
func SetExt(path string, ext string, oldexts ...string) string {
	if ext == `` {
		return path
	}

	var oldext string

	if len(oldexts) > 0 && oldexts[0] != `` {
		oldext = oldexts[0]
	} else {
		oldext = filepath.Ext(path)
	}

	oldext = strings.TrimPrefix(oldext, `.`)
	ext = strings.TrimPrefix(ext, `.`)

	if strings.HasSuffix(path, `.`+oldext) {
		path = strings.TrimSuffix(path, `.`+oldext) + `.` + ext
	}

	return path
}

// Returns the size of the given filename, or zero if any error is encountered (or if the file is zero length).
func SizeOf(path string) convutil.Bytes {
	path = MustExpandUser(path)

	if stat, err := os.Stat(path); err == nil {
		return convutil.Bytes(stat.Size())
	}

	return 0
}

// Compare the binary contents of two io.Reader instances.  The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
// An optional hash.Hash instance may be given, otherwise the default crypto/sha256 will be used.
func CompareReaders(a io.Reader, b io.Reader, hasher ...hash.Hash) int {
	var h hash.Hash

	if len(hasher) > 0 && hasher[0] != nil {
		h = hasher[0]
	} else {
		h = sha256.New()
	}

	if _, err := io.Copy(h, a); err == nil {
		s1 := h.Sum(nil)
		h.Reset()

		if _, err := io.Copy(h, b); err == nil {
			s2 := h.Sum(nil)

			return bytes.Compare(s1, s2)
		}
	}

	return -1
}

// Returns whether two files represent the same file.  If a string is given for either file, os.Stat will be run on that
// path.  If os.FileInfo is given for either file, it will be passed to os.SameFile directly.  If either file is an
// io.Reader, the contents of both files will be read and hashed using CompareHasher.  If the hashes are identical, the
// files are considered the same.  Any error encountered and this function will return false.
func SameFile(first interface{}, second interface{}) bool {
	if first == second {
		return true
	}

	var i1, i2 os.FileInfo
	var r1, r2 io.Reader

	if r, ok := first.(io.Reader); ok {
		r1 = r
	} else if i, ok := first.(os.FileInfo); ok {
		i1 = i
	} else if stat, err := os.Stat(MustExpandUser(typeutil.String(first))); err == nil {
		i1 = stat
	} else {
		return false
	}

	if r, ok := second.(io.Reader); ok {
		r2 = r
	} else if i, ok := second.(os.FileInfo); ok {
		i2 = i
	} else if stat, err := os.Stat(MustExpandUser(typeutil.String(second))); err == nil {
		i2 = stat
	} else {
		return false
	}

	if r1 == nil && r2 == nil {
		return os.SameFile(i1, i2)
	}

	if r1 == nil {
		if f, err := os.Open(i1.Name()); err == nil {
			defer f.Close()
			r1 = f
		} else {
			return false
		}
	}

	if r2 == nil {
		if f, err := os.Open(i2.Name()); err == nil {
			defer f.Close()
			r2 = f
		} else {
			return false
		}
	}

	return (CompareReaders(r1, r2, CompareHasher) == 0)
}
