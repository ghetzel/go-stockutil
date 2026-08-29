package httputil

import (
	"github.com/jdxcode/netrc"
	"go.gary.cool/go-stockutil/fileutil"
	"go.gary.cool/go-stockutil/sliceutil"
)

var NetrcFile = `~/.netrc`

// Retreive the plaintext username and password from the netrc-formatted file in the
// NetrcFile package variable.  The final return argument will be true if and only if
// the .netrc file exists, is readable, and the username OR password matched to the given
// domain is non-empty.
func NetrcCredentials(domain string) (string, string, bool) {
	if domain != `` {
		if path, err := fileutil.ExpandUser(NetrcFile); err == nil {
			if nrc, err := netrc.Parse(path); err == nil {
				if m := nrc.Machine(domain); m != nil {
					var user = sliceutil.OrString(m.Get(`login`), m.Get(`username`))
					var pass = m.Get(`password`)

					if user != `` || pass != `` {
						return user, pass, true
					}
				} else if domain != `*` {
					return NetrcCredentials(`*`)
				}
			}
		}
	}

	return ``, ``, false
}
