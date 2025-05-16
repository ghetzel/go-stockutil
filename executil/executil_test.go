package executil

import (
	"os/exec"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestWhich(t *testing.T) {
	var assert = require.New(t)

	// assert.Equal(`/bin/sh`, Which(`/bin/sh`))
	// assert.Equal(`/bin/sh`, Which(`sh`))
	// assert.True(strings.HasSuffix(Which(`tail`), `/bin/tail`))
	assert.Empty(Which(`absolutely-not-a-command-@#$%^&*`))
}

func TestJoin(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(``, Join(nil))
	assert.Equal(``, Join(``))
	assert.Equal(``, Join([]string{}))
	assert.Equal(`ls -l`, Join([]string{`ls`, `-l`}))
	assert.Equal(`ls -l '/this is a folder'`, Join([]string{`ls`, `-l`, `/this is a folder`}))

	assert.Equal(`ls -l '/this is a folder'`, Join(
		exec.Command(`ls`, `-l`, `/this is a folder`),
	))

	assert.Equal(`whoami`, Join(
		exec.Command(`whoami`),
	))

	assert.Equal(``, Join(
		exec.Command(``),
	))
}

func TestEventedLineOutput(t *testing.T) {
	var assert = require.New(t)

	var stdout = make([]string, 0)
	var stderr = make([]string, 0)

	var cmd = Command(`printf`, `1\n2\n3\n`)
	cmd.OnStdout = func(line string, serr bool) {
		stdout = append(stdout, line)
	}

	assert.NoError(cmd.Run())
	assert.Equal([]string{
		`1`, `2`, `3`,
	}, stdout)
	assert.Empty(stderr)

	// test stdout/stderr interleaving
	cmd = Command(`bash`, `-c`, `echo mock; echo yeah 1>&2; echo ing; echo yeah 1>&2; echo bird; echo yeah 1>&2; echo yeah; echo yeah 1>&2`)
	stdout = nil
	stderr = nil

	cmd.OnStdout = func(line string, serr bool) {
		stdout = append(stdout, line)
	}

	cmd.OnStderr = func(line string, serr bool) {
		stderr = append(stderr, line)
	}

	assert.NoError(cmd.Run())
	assert.Equal([]string{
		`mock`, `ing`, `bird`, `yeah`,
	}, stdout)

	assert.Equal([]string{
		`yeah`, `yeah`, `yeah`, `yeah`,
	}, stderr)
}
