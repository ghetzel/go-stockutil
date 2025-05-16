package fileutil

import (
	"context"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ghetzel/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestRetrieveViaSSH_SFTP(t *testing.T) {
	// example: TEST_GOSTOCKUTIL_RETRIEVE_VIA_SFTP="sftp://localhost:2020/file-in-homedir|sftp://ubuntu:password@remote-host//etc/fstab"
	var sshUrlEnv = os.Getenv(`TEST_GOSTOCKUTIL_RETRIEVE_VIA_SFTP`)

	if sshUrlEnv != `` {
		for _, uri := range strings.Split(sshUrlEnv, `|`) {
			uri = strings.TrimSpace(uri)

			if uri == `` {
				continue
			}

			u, err := url.Parse(uri)

			assert.NoError(t, err)
			assert.NotNil(t, u)

			ctx := context.WithValue(context.Background(), `verifyHostFunc`, func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				var p = u.Port()

				if p == `` {
					p = `22`
				}

				assert.Equal(t, u.Host+`:`+p, hostname)
				return nil
			})

			d, err := RetrieveViaSSH(ctx, u)

			assert.NoError(t, err)

			data, err := io.ReadAll(d)
			assert.NoError(t, err)
			assert.NoError(t, d.Close())
			assert.Equal(t, "HELLO THERE\n", string(data))
		}
	}
}

func TestRetrieveViaSSH_SSH(t *testing.T) {
	// example: TEST_GOSTOCKUTIL_RETRIEVE_VIA_SSH="ssh://localhost:2020/hostname"
	var sshUrlEnv = os.Getenv(`TEST_GOSTOCKUTIL_RETRIEVE_VIA_SSH`)

	if sshUrlEnv != `` {
		for _, uri := range strings.Split(sshUrlEnv, `|`) {
			uri = strings.TrimSpace(uri)

			if uri == `` {
				continue
			}

			u, err := url.Parse(uri)

			assert.NoError(t, err)
			assert.NotNil(t, u)

			ctx := context.WithValue(context.Background(), `verifyHostFunc`, func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				var p = u.Port()

				if p == `` {
					p = `22`
				}

				assert.Equal(t, u.Host+`:`+p, hostname)
				return nil
			})

			d, err := RetrieveViaSSH(ctx, u)

			assert.NoError(t, err)

			data, err := io.ReadAll(d)

			assert.NoError(t, err)
			assert.NoError(t, d.Close())
			assert.Equal(t, u.Hostname()+"\n", string(data))
		}
	}
}
