package executil

import (
	"bytes"
	"io"
	"testing"

	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/testify/assert"
)

func incr(value *int) CommandStatusFunc {
	return func(s Status) {
		if value != nil {
			*value += 1
		}
	}
}

func TestExecTrue(t *testing.T) {
	var starts int
	var completes int
	var successes int
	var failures int
	var cmd = Command(`true`)

	cmd.OnStart = incr(&starts)
	cmd.OnComplete = incr(&completes)
	cmd.OnSuccess = incr(&successes)
	cmd.OnError = incr(&failures)

	assert.False(t, cmd.Status().Running)
	assert.False(t, cmd.Status().Successful)
	assert.Zero(t, cmd.Status().PID)
	assert.True(t, cmd.Status().StartedAt.IsZero())
	assert.True(t, cmd.Status().StoppedAt.IsZero())
	assert.Zero(t, cmd.Status().ExitCode)
	assert.Nil(t, cmd.Status().Error)

	cmd.Run()

	assert.False(t, cmd.Status().Running)
	assert.True(t, cmd.Status().Successful)
	assert.False(t, cmd.Status().PID == 0)
	assert.False(t, cmd.Status().StartedAt.IsZero())
	assert.False(t, cmd.Status().StoppedAt.IsZero())
	assert.True(t, cmd.Status().Took() > 0)
	assert.Zero(t, cmd.Status().ExitCode)
	assert.Nil(t, cmd.Status().Error)

	assert.Equal(t, 1, starts)
	assert.Equal(t, 1, completes)
	assert.Equal(t, 1, successes)
	assert.Equal(t, 0, failures)
}

func TestExecFalse(t *testing.T) {
	var starts int
	var completes int
	var successes int
	var failures int
	var cmd = Command(`false`)

	cmd.OnStart = incr(&starts)
	cmd.OnComplete = incr(&completes)
	cmd.OnSuccess = incr(&successes)
	cmd.OnError = incr(&failures)

	assert.False(t, cmd.Status().Running)
	assert.False(t, cmd.Status().Successful)
	assert.Zero(t, cmd.Status().PID)
	assert.True(t, cmd.Status().StartedAt.IsZero())
	assert.True(t, cmd.Status().StoppedAt.IsZero())
	assert.Zero(t, cmd.Status().ExitCode)
	assert.Nil(t, cmd.Status().Error)

	cmd.Run()

	assert.False(t, cmd.Status().Running)
	assert.False(t, cmd.Status().Successful)
	assert.False(t, cmd.Status().PID == 0)
	assert.False(t, cmd.Status().StartedAt.IsZero())
	assert.False(t, cmd.Status().StoppedAt.IsZero())
	assert.True(t, cmd.Status().Took() > 0)
	assert.True(t, cmd.Status().ExitCode == 1)
	assert.EqualError(t, cmd.Status().Error, `Process exited with status 1`)

	assert.Equal(t, 1, starts)
	assert.Equal(t, 1, completes)
	assert.Equal(t, 0, successes)
	assert.Equal(t, 1, failures)
}

func TestShellOut(t *testing.T) {
	assert.Equal(t, `hello there`, string(MustShellOut(`echo`, `-n`, `hello`, `there`)))
	assert.Equal(t, `hello there`, string(MustShellOut(`echo -n hello there`)))
	assert.Equal(t, `hello there`, string(MustShellOut(`echo -n`, `hello there`)))
}

func TestExecReadCloser(t *testing.T) {
	var cmd = Command(`echo`, `hello`)
	var data, err = io.ReadAll(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "hello\n", string(data))
	assert.NoError(t, cmd.Close())
}

func TestExecWriteCloser(t *testing.T) {
	var cmd = Command(`cat`)

	var n, err = cmd.Write([]byte("hello\n"))
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	assert.NoError(t, cmd.Close())
}

func TestExecReadWriteCloser(t *testing.T) {
	var c int = 128
	var payload = stringutil.UUID().Bytes()
	var cmd = Command(`cat`)

	for i := 0; i < c; i++ {
		var _, werr = cmd.Write(payload)
		assert.NoError(t, werr)
	}

	assert.NoError(t, cmd.CloseInput())

	var data, rerr = io.ReadAll(cmd)
	assert.NoError(t, rerr)
	assert.Equal(t, bytes.Repeat(payload, c), data)
	// assert.NoError(t, cmd.Close())
}
