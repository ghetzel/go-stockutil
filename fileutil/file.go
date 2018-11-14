package fileutil

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghetzel/go-stockutil/pathutil"
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

func IsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func ReadAll(filename string) ([]byte, error) {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		return ioutil.ReadAll(file)
	} else {
		return nil, err
	}
}

func ReadAllString(filename string) (string, error) {
	if data, err := ReadAll(filename); err == nil {
		return string(data), nil
	} else {
		return ``, err
	}
}

func ReadAllLines(filename string) ([]string, error) {
	if data, err := ReadAllString(filename); err == nil {
		return strings.Split(data, "\n"), nil
	} else {
		return nil, err
	}
}

func ReadFirstLine(filename string) (string, error) {
	if lines, err := ReadAllLines(filename); err == nil {
		return lines[0], nil
	} else {
		return ``, err
	}
}

func MustReadAll(filename string) []byte {
	if data, err := ReadAll(filename); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}

func MustReadAllString(filename string) string {
	if data, err := ReadAllString(filename); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}
