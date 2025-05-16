package fileutil

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ghetzel/go-stockutil/pathutil"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const DefaultOpenTimeout = time.Duration(10 * time.Second)

type OpenOptions struct {
	Timeout  time.Duration
	Insecure bool
}

type OpenHandler func(*url.URL, OpenOptions) (io.ReadCloser, error)

func (opts OpenOptions) GetTimeout() time.Duration {
	if opts.Timeout == 0 {
		return DefaultOpenTimeout
	} else {
		return opts.Timeout
	}
}

var openHandlers = map[string]OpenHandler{
	``:      openHandlerLocalFile,
	`file`:  openHandlerLocalFile,
	`http`:  openHandlerHttp,
	`https`: openHandlerHttp,
	`ftp`:   openHandlerFtp,
	`sftp`:  openHandlerSftp,
}

func openHandlerLocalFile(uri *url.URL, opt OpenOptions) (io.ReadCloser, error) {
	if expanded, err := ExpandUser(uri.Path); err == nil {
		return os.Open(expanded)
	} else {
		return nil, err
	}
}

func openHandlerHttp(uri *url.URL, opt OpenOptions) (io.ReadCloser, error) {
	var client = http.Client{
		Timeout: opt.GetTimeout(),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: opt.Insecure,
			},
		},
	}

	if response, err := client.Get(uri.String()); err == nil {
		if response != nil {
			return response.Body, nil
		} else {
			uri.User = nil
			return nil, fmt.Errorf("empty response from %v", uri.String())
		}
	} else if response != nil {
		uri.User = nil
		return nil, fmt.Errorf("from %v: %v", uri.String(), response.Status)
	} else {
		return nil, fmt.Errorf("empty response from %v", uri.String())
	}

}

func openHandlerFtp(u *url.URL, opt OpenOptions) (io.ReadCloser, error) {
	if u.Port() == `` {
		u.Host += `:21`
	}

	if conn, err := ftp.DialTimeout(u.Host, opt.GetTimeout()); err == nil {
		if up := u.User; up != nil {
			pw, _ := up.Password()

			if err := conn.Login(up.Username(), pw); err != nil {
				u.User = nil
				return nil, fmt.Errorf("from %v: %v", u.String(), err)
			}
		} else if err := conn.Login(`anonymous`, `anonymous@example.com`); err != nil {
			u.User = nil
			return nil, fmt.Errorf("from %v: %v", u.String(), err)
		}

		if response, err := conn.Retr(u.Path); err == nil {
			return NewPostReadCloser(response, func(rc io.ReadCloser) error {
				defer conn.Quit()
				return rc.Close()
			}), nil
		} else {
			defer conn.Quit()
			return nil, err
		}
	} else {
		u.User = nil
		return nil, fmt.Errorf("from %v: %v", u.String(), err)
	}
}

func openHandlerSftp(u *url.URL, opts OpenOptions) (io.ReadCloser, error) {
	var username string
	var methods []ssh.AuthMethod

	if u.Port() == `` {
		u.Host += `:22`
	}

	if user := u.User; user != nil {
		username = user.Username()

		if pw, ok := user.Password(); ok {
			methods = append(methods, ssh.Password(pw))
		}
	}

	if username == `` {
		username = os.Getenv(`USER`)
	}

	// go through the whole process of loading keypairs
	if v := u.Query().Get(`keyfile`); v != `` {
		if kexpanded, err := pathutil.ExpandUser(v); err == nil {
			if key, err := os.ReadFile(kexpanded); err == nil {
				if signer, err := ssh.ParsePrivateKey(key); err == nil {
					methods = append(methods, ssh.PublicKeys(signer))
				} else {
					return nil, fmt.Errorf("unable to parse private key: %v", err)
				}
			} else {
				return nil, fmt.Errorf("unable to load private key: %v", err)
			}
		} else {
			return nil, err
		}
	}

	// by this point we should have a username and at least one method
	// authenticating in the pipe
	if len(methods) == 0 {
		return nil, fmt.Errorf("no SSH authentication methods specified; provide a password or ?keyfile query string")
	}

	if username != `` {
		if client, err := ssh.Dial(`tcp`, u.Host, &ssh.ClientConfig{
			User: username,
			Auth: methods,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}); err == nil {
			if sclient, err := sftp.NewClient(client); err == nil {
				if file, err := sclient.Open(u.Path); err == nil {
					return NewPostReadCloser(file, func(rc io.ReadCloser) error {
						defer sclient.Close()
						return rc.Close()
					}), nil
				} else {
					return nil, err
				}
			} else {
				u.User = nil
				return nil, fmt.Errorf("from %v: %v", u.String(), err)
			}
		} else {
			u.User = nil
			return nil, fmt.Errorf("from %v: %v", u.String(), err)
		}
	} else {
		return nil, fmt.Errorf("%v: must specify a username", u.Scheme)
	}
}

// Register a handler for a new or existing URL scheme, for use with Open() and OpenWithOptions()
func RegisterOpenHandler(scheme string, handler OpenHandler) {
	openHandlers[scheme] = handler
}

// Removes a registered URL scheme handler.
func RemoveOpenHandler(scheme string) {
	delete(openHandlers, scheme)
}

// Calls OpenWithOptions with no options set.
func Open(uri string) (io.Reader, error) {
	return OpenWithOptions(uri, OpenOptions{})
}

// A generic URL opener that supports various schemes and returns an io.Reader.
// Supported URL schemes include: file://, http://, https://, ftp://, sftp://. If no scheme is
// provided, the URL is interpreted as a local filesystem path.
func OpenWithOptions(uri string, options OpenOptions) (io.Reader, error) {
	if u, err := url.Parse(uri); err == nil {
		if handler, ok := openHandlers[u.Scheme]; ok {
			return handler(u, options)
		} else {
			return nil, fmt.Errorf("unsupported scheme %s", u.Scheme)
		}
	} else {
		return nil, fmt.Errorf("invalid source URL or filename")
	}
}
