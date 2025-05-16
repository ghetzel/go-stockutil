package fileutil

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"os/user"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/mattn/go-shellwords"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshHostKeyCallbackFunc = func(hostname string, remote net.Addr, key ssh.PublicKey) error

var SshPrivateKey, _ = ExpandUser(`~/.ssh/id_rsa`)
var SshVerifyHostFunc ssh.HostKeyCallback
var SshDefaultTimeout = 10 * time.Second

// Retrieve a file via SFTP (SSH file transfer).  The given URL should resemble the
// prototype: ssh://[user:password@]hostname[:22]/path/relative/to/homedir
//
// This function will honor any authentication details from a running SSH agent,
// as well as utilize the private key located in the path indicated by SshPrivateKey
// or via the `privateKey` context value.
//
// Supported Context Values:
//
// username:
//
//	(string) the username to login with. can be overriden by a username specified in the URL.
//
// password:
//
//	(string) the password to login with. can be overriden by a password specified in the URL.
//
// passphrase:
//
//	(string) context value specifies a plaintext passphrase used to unlock the local private keyfile.
//
// insecure:
//
//	(bool) whether to ignore remote hostkey checks.  Does not work if verifyHostFunc is set.
//
// verifyHostFunc:
//
//	(SshHostKeyCallbackFunc) context value, if it is convertible to the ssh.HostKeyCallback type, will
//	be called to verify the remote SSH host key in a manner of the function's choosing.  The default
//	behavior is to accept all remote hostkeys as valid.
func RetrieveViaSSH(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	ctx, timeout := ctxToTimeout(ctx, SshDefaultTimeout)

	var authMethods goph.Auth
	var username string = typeutil.String(ctx.Value(`username`))
	var password string = typeutil.String(ctx.Value(`password`))
	var port int = typeutil.OrNInt(u.Port(), 22)
	var remotePath = strings.TrimPrefix(u.Path, `/`)
	var keyPassphrase = typeutil.String(ctx.Value(`passphrase`))
	var verifyHostFunc = SshVerifyHostFunc

	if vhfn, ok := ctx.Value(`verifyHostFunc`).(SshHostKeyCallbackFunc); ok {
		verifyHostFunc = vhfn
	} else if typeutil.Bool(ctx.Value(`insecure`)) {
		verifyHostFunc = ssh.InsecureIgnoreHostKey()
	}

	var keyFile = sliceutil.OrString(
		typeutil.String(ctx.Value(`passphrase`)),
		SshPrivateKey,
	)

	if a, err := goph.UseAgent(); err == nil {
		authMethods = append(authMethods, a...)
	}

	if ui := u.User; ui != nil {
		if u := ui.Username(); u != `` {
			username = u
		}

		if p, ok := ui.Password(); ok {
			password = p
		}
	}

	if password != `` {
		authMethods = append(authMethods, goph.Password(password)...)
	}

	if IsNonemptyFile(keyFile) {
		if a, err := goph.Key(keyFile, keyPassphrase); err == nil {
			authMethods = append(authMethods, a...)
		}
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no client authentication methods available")
	}

	if username == `` {
		if cur, err := user.Current(); err == nil {
			username = cur.Username
		} else {
			return nil, err
		}
	}

	if client, err := goph.NewConn(&goph.Config{
		User:     username,
		Addr:     u.Hostname(),
		Port:     uint(port),
		Auth:     authMethods,
		Timeout:  timeout,
		Callback: verifyHostFunc,
	}); err == nil {
		var readCloser io.ReadCloser
		var rerr error

		switch u.Scheme {
		case `sftp`:
			if sftp, err := sftp.NewClient(client.Client); err == nil {
				if file, err := sftp.Open(remotePath); err == nil {
					// setup a post-read closer that will handle closing both the remote file handle
					// and the remote connection (in that order)
					readCloser = NewPostReadCloser(file, func(rc io.ReadCloser) error {
						defer sftp.Close()
						defer client.Close()

						return file.Close()
					})
				} else {
					defer sftp.Close()
					return nil, err
				}
			} else {

				return nil, err
			}
		case `ssh`:
			if argv, err := shellwords.Parse(remotePath); err == nil {
				if scmd, err := client.Command(argv[0], argv[1:]...); err == nil {
					if out, err := scmd.StdoutPipe(); err == nil {
						if err := scmd.Start(); err == nil {
							readCloser = NewPostReadCloser(io.NopCloser(out), func(rc io.ReadCloser) error {
								defer client.Close()
								return scmd.Wait()
							})
						} else {
							rerr = fmt.Errorf("RetrieveViaSSH: command failed to start: %v", err)
						}
					} else {
						rerr = fmt.Errorf("RetrieveViaSSH: bad pipe: %v", err)
					}
				} else {
					rerr = fmt.Errorf("RetrieveViaSSH: bad shell command: %v", err)
				}
			} else {
				rerr = fmt.Errorf("RetrieveViaSSH: bad shell command: %v", err)
			}

		default:
			rerr = fmt.Errorf("RetrieveViaSSH: bad scheme %q", u.Scheme)
		}

		// if the readcloser is non-nil, then it will close the connection itself after
		// the read is completed (with Close()).  if it is nil, then we have to close the
		// connection ourselves
		if readCloser == nil {
			defer client.Close()
		}

		return readCloser, rerr
	} else {
		return nil, err
	}
}
