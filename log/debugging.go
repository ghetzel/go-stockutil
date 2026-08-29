package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gobwas/glob"
	"go.gary.cool/go-stockutil/stringutil"
)

var MaxStackTraceDepth = 32

type StackItem struct {
	ProgramCounter uintptr
	Filename       string
	Line           int
	Function       string
	PackageName    string
	Receiver       string
	FunctionName   string
}

func (self StackItem) InPackage(pkgname string) bool {
	if self.PackageName == pkgname {
		return true
	} else if glob.MustCompile(pkgname).Match(self.PackageName) {
		return true
	}

	return false
}

func (self StackItem) String() string {
	var line []string

	if self.Function != `` {
		line = append(line, fmt.Sprintf("function %v", self.Function))
	}

	if self.Filename != `` && self.Line > 0 {
		line = append(line, fmt.Sprintf("%s, line %d", self.Filename, self.Line))
	}

	return strings.Join(line, "\n")
}

type StackItems []StackItem

func (self StackItems) String() string {
	var lines []string

	for _, item := range self {
		lines = append(lines, item.String(), "")
	}

	return strings.Join(lines, "\n")
}

// Retrieves details about the call stack that led to this function call.
func StackTrace(skip int) StackItems {
	var pc = make([]uintptr, MaxStackTraceDepth)
	var items = make(StackItems, 0)

	if n := runtime.Callers(skip, pc); n <= len(pc) {
		pc = pc[:n]

		if frames := runtime.CallersFrames(pc); frames != nil {
			for i := 0; len(items) <= len(pc); i++ {
				var frame, more = frames.Next()
				var pkgname string
				var recv string
				var fnname string
				var lastSlash = strings.LastIndex(frame.Function, `/`)

				if lastSlash > 0 {
					var pkgRecvSep = lastSlash + strings.Index(frame.Function[lastSlash:], `.`)
					var recvFnnSep = strings.LastIndex(frame.Function, `.`)

					if recvFnnSep > 0 && pkgRecvSep > 0 {
						pkgname = frame.Function[0:pkgRecvSep]

						if recvFnnSep != pkgRecvSep {
							recv = frame.Function[pkgRecvSep+1 : recvFnnSep]
						}

						fnname = frame.Function[recvFnnSep+1:]
					}
				} else if lastSlash < 0 {
					pkgname, fnname = stringutil.SplitPair(frame.Function, `.`)
				}

				items = append(items, StackItem{
					ProgramCounter: frame.PC,
					Function:       frame.Function,
					PackageName:    pkgname,
					Receiver:       recv,
					FunctionName:   fnname,
					Filename:       frame.File,
					Line:           frame.Line,
				})

				if !more {
					break
				}
			}
		}
	}

	return items
}

// Logs the current stack trace as debug log output.
func DebugStack() {
	Debug("Stack trace:")

	for i, item := range StackTrace(3) {
		for j, line := range strings.Split(item.String(), "\n") {
			if j == 0 {
				Debugf("  % 2d: %v", i, line)
			} else {
				Debugf("          %v", line)
			}
		}
	}
}
