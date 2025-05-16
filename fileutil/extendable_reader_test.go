package fileutil

import (
	"bytes"
	"io"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestExtendableReader(t *testing.T) {
	var seqrc ExtendableReader
	var b = make([]byte, 64)

	var n, err = seqrc.Read(b)
	require.Zero(t, n)
	require.Equal(t, io.EOF, err)

	seqrc.AppendSource(io.NopCloser(bytes.NewBufferString(`hello`)))

	// should read exactly 5 bytes
	n, err = seqrc.Read(b)
	require.NoError(t, err)
	require.Equal(t, 5, n)

	// should EOF
	n, err = seqrc.Read(b)
	require.Zero(t, n)
	require.Equal(t, io.EOF, err)

	seqrc.AppendSource(io.NopCloser(bytes.NewBufferString(`there`)))
	seqrc.AppendSource(io.NopCloser(bytes.NewBufferString(` today`)))
	seqrc.AppendSource(io.NopCloser(bytes.NewBufferString(` is`)))
	seqrc.AppendSource(io.NopCloser(bytes.NewBufferString(` a day`)))

	b, err = io.ReadAll(&seqrc)
	require.NoError(t, err)
	require.Equal(t, `there today is a day`, string(b))

	b, err = io.ReadAll(&seqrc)
	require.NoError(t, err)
	require.Len(t, b, 0)
}
