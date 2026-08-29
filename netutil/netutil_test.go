package netutil

// import (
// 	"os/exec"
// 	"strings"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func TestFQDN(t *testing.T) {
// 	assert := require.New(t)
// 	sys, err := exec.Command(`hostname`, `-f`).Output()
// 	assert.NoError(err)

// 	var want = strings.TrimSpace(string(sys))
// 	assert.Equal(want, FQDN())
// }

// func TestPingLocalhost4(t *testing.T) {
// 	assert := require.New(t)
// 	assert.NoError(OnePingOnly(net.ParseIP(`127.0.0.1`), nil, 10*time.Second))

// 	defIP := DefaultAddress()
// 	assert.NoError(OnePingOnly(defIP.Gateway, defIP, 10*time.Second))
// }

// func TestPingLocalhost6(t *testing.T) {
// 	assert := require.New(t)
// 	assert.NoError(OnePingOnly(net.ParseIP(`::1`), nil, 10*time.Second))

// 	defIP := DefaultAddress()
// 	assert.NoError(OnePingOnly(defIP.Gateway, nil, 10*time.Second))
// }

// func TestWaitForPing4(t *testing.T) {
// 	assert := require.New(t)
// 	assert.NoError(WaitForPing(`127.0.0.1`, 10*time.Second))
// 	assert.NoError(WaitForPing(net.ParseIP(`127.0.0.1`), 10*time.Second))
// 	assert.NoError(WaitForPing(&IPAddress{
// 		IP: net.ParseIP(`127.0.0.1`),
// 	}, 10*time.Second))
// }

// func TestWaitForPing6(t *testing.T) {
// 	assert := require.New(t)
// 	assert.NoError(WaitForPing(`::1`, 10*time.Second))
// 	assert.NoError(WaitForPing(net.ParseIP(`::1`), 10*time.Second))
// 	assert.NoError(WaitForPing(&IPAddress{
// 		IP: net.ParseIP(`::1`),
// 	}, 10*time.Second))
// }

// func TestWaitForGatewayPing4(t *testing.T) {
// 	assert := require.New(t)
// 	assert.NoError(WaitForGatewayPing(10 * time.Second))
// }

// func TestWaitForGatewayPing6(t *testing.T) {
// 	assert := require.New(t)
// 	assert.NoError(WaitForGatewayPing6(10 * time.Second))
// }
