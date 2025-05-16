package fileutil

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/ghetzel/testify/require"
)

const testTextPre = "// this is a cool test\n"
const testTextBody = ("\n" +
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut \n" +
	"labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco \n" +
	"laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in \n" +
	"voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat \n" +
	"cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." +
	"\n" +
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut \n" +
	"labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco \n" +
	"laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in \n" +
	"voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat \n" +
	"cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

const testText = testTextPre + testTextBody

func tt() io.Reader {
	return bytes.NewBufferString(testText)
}

func TestReadManipulatorNoOp(t *testing.T) {
	var assert = require.New(t)

	var rm = NewReadManipulator(tt())
	out, err := io.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(testText, string(out))
}

func TestReadManipulatorDoNothingFunction(t *testing.T) {
	var assert = require.New(t)

	var rm = NewReadManipulator(tt(), func(in []byte) ([]byte, error) {
		return in, nil
	})

	out, err := io.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(testText, string(out))
}

func TestReadManipulatorRemoveComments(t *testing.T) {
	var assert = require.New(t)

	var rm = NewReadManipulator(tt(), RemoveLinesWithPrefix(`//`, true))

	out, err := io.ReadAll(rm)
	assert.NoError(err)
	assert.Equal("\n"+testTextBody, string(out))
}

func TestReadManipulatorRemoveBlankLines(t *testing.T) {
	var assert = require.New(t)

	var rm = NewReadManipulator(tt(), RemoveBlankLines)

	out, err := io.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(strings.TrimSpace(testTextPre)+testTextBody, string(out))
}

func TestReadManipulatorDolorToBacon(t *testing.T) {
	var assert = require.New(t)

	var rm = NewReadManipulator(tt(), func(in []byte) ([]byte, error) {
		var line = strings.Replace(string(in), `dolor`, `bacon`, -1)
		return []byte(line), nil
	})

	out, err := io.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(
		strings.Replace(testText, `dolor`, `bacon`, -1),
		string(out),
	)
}

func TestReadManipulatorManipulateAll(t *testing.T) {
	var assert = require.New(t)

	var fn1 = ReplaceWith(`dolor`, `bacon`, -1)
	var fn2 = RemoveBlankLines
	var fn3 = func(in []byte) ([]byte, error) {
		var line = strings.Replace(string(in), `nostrud`, `potato`, -1)
		return []byte(line), nil
	}

	var rm = NewReadManipulator(tt(), ManipulateAll(fn1, fn2, fn3))

	var wanted = testText
	wanted = strings.Replace(wanted, `dolor`, `bacon`, -1)
	wanted = strings.Replace(wanted, `nostrud`, `potato`, -1)
	wanted = strings.Replace(wanted, "\n\n", "\n", -1)

	out, err := io.ReadAll(rm)
	assert.NoError(err)
	assert.Equal(wanted, string(out))
}
