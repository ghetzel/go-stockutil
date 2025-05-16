package executil

import (
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/mattn/go-shellwords"
)

// Locates the first path containing the given command. The directories listed
// in the environment variable "PATH" will be checked, in order.  If additional
// directories are specified in the path variadic argument, they will be checked
// first.  If the command is not in any path, an empty string will be returned.
func Which(cmdname string, path ...string) string {
	if found := WhichAll(cmdname, path...); len(found) > 0 {
		return found[0]
	} else {
		return ``
	}
}

// Locates the all paths containing the given command. The directories listed
// in the environment variable "PATH" will be checked, in order.  If additional
// directories are specified in the path variadic argument, they will be checked
// first.  If the command is not in any path, an empty slice will be returned.
func WhichAll(cmdname string, path ...string) []string {
	var dirs = append(path, strings.Split(os.Getenv(`PATH`), `:`)...)
	var found = make([]string, 0)

	if fileutil.IsNonemptyExecutableFile(cmdname) {
		found = append(found, cmdname)
	}

	for _, dir := range dirs {
		var candidate = filepath.Join(dir, cmdname)

		if len(strings.TrimSpace(dir)) == 0 {
			continue
		} else if !fileutil.DirExists(dir) {
			continue
		} else if fileutil.IsNonemptyExecutableFile(candidate) {
			found = append(found, candidate)
		}
	}

	return found
}

// Take an *exec.Cmd or []string and return a shell-executable command line string.
func Join(in any) string {
	var args []string

	if cmd, ok := in.(*exec.Cmd); ok {
		args = cmd.Args
	} else if typeutil.IsArray(in) {
		args = sliceutil.Stringify(in)
	} else {
		return ``
	}

	for i, arg := range args {
		if strings.Contains(arg, ` `) {
			args[i] = stringutil.Wrap(arg, `'`, `'`)
		}
	}

	return strings.Join(args, ` `)
}

// Uses environment variables and other configurations to attempt to locate the
// path to the user's shell.
func FindShell() string {
	var shells = []string{os.Getenv(`SHELL`)}
	shells = append(shells, Which(`sh`))

	for _, shell := range shells {
		if shell != `` {
			return shell
		}
	}

	return ``
}

// Returns whether the current user is root or not.
func IsRoot() bool {
	if current, err := user.Current(); err == nil {
		if current.Uid == `0` {
			return true
		}
	}

	return false
}

// Returns the first argument if the current user is root, and the second if not.
func RootOr(ifRoot any, notRoot any) any {
	if IsRoot() {
		return ifRoot
	} else {
		return notRoot
	}
}

// The same as RootOr, but returns a string.
func RootOrString(ifRoot any, notRoot any) string {
	if IsRoot() {
		return typeutil.String(ifRoot)
	} else {
		return typeutil.String(notRoot)
	}
}

// Registers a list of OS signals to intercept and provides an opportunity to run
// a function before the program exits.
func TrapSignals(after func(sig os.Signal) bool, signals ...os.Signal) {
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	for trap := range signalChan {
		if after != nil {
			if !after(trap) {
				return
			}
		}
	}
}

// Splits the given string into words, honoring quoting and escaping conventions of common command line shells.
func Split(cmd string) ([]string, error) {
	return shellwords.Parse(cmd)
}

// Same as Split, but silently discards any errors, returning an empty slice in this case.
func ShouldSplit(cmd string) []string {
	if words, err := Split(cmd); err == nil {
		return words
	} else {
		return nil
	}
}

// Same as Split, but panics if there is an error.
func MustSplit(cmd string) []string {
	if words, err := Split(cmd); err == nil {
		return words
	} else {
		panic(err.Error())
	}
}
