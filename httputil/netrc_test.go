package httputil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.gary.cool/go-stockutil/fileutil"
)

func TestNetrcPositive(t *testing.T) {
	NetrcFile = fileutil.MustWriteTempFile(
		"machine stock-httputil-test login hello password there\nmachine intentionally-left-blank\n",
		"test-ghetzel-go-stockutil-httputil",
	)

	defer os.Remove(NetrcFile)

	var u, p, ok = NetrcCredentials(``)

	require.False(t, ok)
	require.Empty(t, u)
	require.Empty(t, p)

	u, p, ok = NetrcCredentials(`nope`)

	require.False(t, ok)
	require.Empty(t, u)
	require.Empty(t, p)

	u, p, ok = NetrcCredentials(`intentionally-left-blank`)

	require.False(t, ok)
	require.Empty(t, u)
	require.Empty(t, p)

	u, p, ok = NetrcCredentials(`stock-httputil-test`)

	require.True(t, ok)
	require.Equal(t, `hello`, u)
	require.Equal(t, `there`, p)
}
