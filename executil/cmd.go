// Utilities that make executing commands on the local system a little bit easier.
package executil

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/mattn/go-shellwords"
)

type CommandStatusFunc func(Status)
type OutputLineFunc func(string, bool)

type Status struct {
	StartedAt  time.Time
	StoppedAt  time.Time
	Running    bool
	Successful bool
	ExitCode   int
	Error      error
	PID        int
	Cmd        *Cmd
}

func (self Status) Took() time.Duration {
	if !self.StartedAt.IsZero() {
		if !self.StoppedAt.IsZero() {
			return self.StoppedAt.Sub(self.StartedAt)
		}
	}

	return 0
}

func (self Status) String() string {
	if self.Error != nil {
		return self.Error.Error()
	} else if self.Running {
		if self.PID > 0 {
			return fmt.Sprintf("PID %d has been running for %v", self.PID, time.Since(self.StartedAt))
		} else {
			return fmt.Sprintf("Process has been running for %v, PID unknown", time.Since(self.StartedAt))
		}
	} else if !self.StoppedAt.IsZero() {
		if self.PID > 0 {
			return fmt.Sprintf("PID %d exited with status %d, took %v", self.PID, self.ExitCode, self.Took())
		} else {
			return fmt.Sprintf("Process exited with status %d, took %v, PID unknown", self.ExitCode, self.Took())
		}
	} else if self.StartedAt.IsZero() {
		return fmt.Sprintf("Process has not started yet")
	} else {
		return fmt.Sprintf("Process status is unknown")
	}
}

// A wrapper for exec.Cmd that provides helpful callbacks and monitoring details that are challenging
// to implement.
type Cmd struct {
	*exec.Cmd

	// An interval of time on which the command should be actively checked for run and exit status.
	MonitorInterval time.Duration

	// How long the command may run for before being killed.
	Timeout time.Duration

	// Whether the command invocation should inherit the environment variables of the calling process.
	InheritEnv bool

	// Called when immediately before the command is executed.
	OnStart CommandStatusFunc

	// Called whenever the monitor check is performed.
	OnMonitor CommandStatusFunc

	// Called when the command exits, regardless of success or failure.
	OnComplete CommandStatusFunc

	// Called when the command exits with a non-error status (code 0)
	OnSuccess CommandStatusFunc

	// Called when the command exits with an error status (non-zero exit code, security, invocation, or resource error)
	OnError CommandStatusFunc

	// Called when a line of standard output is written.
	OnStdout OutputLineFunc

	// Called when a line of standard error is written.
	OnStderr OutputLineFunc

	// If specified, this function will determine how to tokenize the stdout stream and when to call OnStdout.  Defaults to bufio.ScanLines.
	StdoutSplitFunc bufio.SplitFunc

	// If specified, this function will determine how to tokenize the stderr stream and when to call OnStderr.  Defaults to bufio.ScanLines.
	StderrSplitFunc bufio.SplitFunc

	// Specifies that the spawned process should inherit the same Process Group ID (PGID) as the parent.
	InheritParent bool

	status      Status
	finalStatus *Status
	statusLock  sync.Mutex
	reallyDone  *sync.WaitGroup
	finished    chan bool
	exitError   error
	reader      io.ReadCloser
	writer      io.WriteCloser
	hasWaited   bool
}

func Wrap(cmd *exec.Cmd) *Cmd {
	return new(cmd)
}

func Command(name string, arg ...string) *Cmd {
	return new(exec.Command(name, arg...))
}

func ShellCommand(cmdline string) *Cmd {
	if shell := FindShell(); shell != `` {
		return Command(shell, `-c`, cmdline)
	}

	return nil
}

// Run a command and return the standard output.  If the first argument contains
// a command and its arguments, it will be executed in the user's shell using FindShell.
// Otherwise, the first argument will be treated as a command and the remaining arguments
// will be passed in parameterized.
func ShellOut(cmdOrLine string, args ...any) ([]byte, error) {
	var cmd *Cmd

	if va, err := shellwords.Parse(cmdOrLine); err == nil {
		if len(va) == 1 {
			cmd = Command(va[0], sliceutil.Stringify(args)...)
		} else {
			cmd = ShellCommand(strings.Join(
				append(va, sliceutil.Stringify(args)...),
				` `,
			))
		}

		if out, err := cmd.Output(); err == nil {
			return out, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// A panicky version of ShellOut.
func MustShellOut(cmdOrLine string, args ...any) []byte {
	if out, err := ShellOut(cmdOrLine, args...); err == nil {
		return out
	} else {
		panic(err.Error())
	}
}

// Attempts to call ShellOut, but will return nil if there is an error.  Does not panic.
func ShouldShellOut(cmdOrLine string, args ...any) []byte {
	if out, err := ShellOut(cmdOrLine, args...); err == nil {
		return out
	} else {
		return nil
	}
}

func new(wrap *exec.Cmd) *Cmd {
	return &Cmd{
		Cmd:             wrap,
		MonitorInterval: (500 * time.Millisecond),
		status:          Status{},
		finished:        make(chan bool),
		reallyDone:      &sync.WaitGroup{},
	}
}

func (self *Cmd) prestart() error {
	self.statusLock.Lock()
	self.hasWaited = false
	self.finalStatus = nil
	self.status = Status{}
	self.status.StartedAt = time.Now()
	self.status.Running = true
	self.statusLock.Unlock()

	if self.InheritParent {
		self.Cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    os.Getpid(),
		}
	}

	self.updateStatus()

	if fn := self.OnStart; fn != nil {
		fn(self.status)
	}

	go self.startMonitoringCommand()

	if self.InheritEnv {
		self.Cmd.Env = append(os.Environ(), self.Cmd.Env...)
	}

	if fn := self.OnStdout; fn != nil {
		if out, err := self.StdoutPipe(); err == nil {
			go func(rc io.ReadCloser) {
				defer rc.Close()

				var splitfn = self.StdoutSplitFunc

				if splitfn == nil {
					splitfn = bufio.ScanLines
				}

				var splitter = stringutil.NewScanInterceptor(splitfn)
				var scanner = bufio.NewScanner(rc)
				scanner.Split(splitter.Scan)

				for scanner.Scan() {
					fn(scanner.Text(), false)
				}
			}(out)
		} else {
			return fmt.Errorf("stdout: %v", err)
		}
	}

	if fn := self.OnStderr; fn != nil {
		if serr, err := self.StderrPipe(); err == nil {
			go func(rc io.ReadCloser) {
				defer rc.Close()

				var splitfn = self.StderrSplitFunc

				if splitfn == nil {
					splitfn = bufio.ScanLines
				}

				var splitter = stringutil.NewScanInterceptor(splitfn)
				var scanner = bufio.NewScanner(rc)
				scanner.Split(splitter.Scan)

				for scanner.Scan() {
					fn(scanner.Text(), false)
				}
			}(serr)
		} else {
			return fmt.Errorf("stdout: %v", err)
		}
	}

	return nil
}

func (self *Cmd) CombinedOutput() ([]byte, error) {
	defer self.killAndWait()

	if err := self.prestart(); err == nil {
		var out, err = self.Cmd.CombinedOutput()
		self.exitError = err
		self.updateStatus()
		return out, err
	} else {
		return nil, err
	}
}

func (self *Cmd) String() string {
	if out, err := self.Output(); err == nil {
		return string(out)
	} else {
		return ``
	}
}

func (self *Cmd) Output() ([]byte, error) {
	defer self.killAndWait()

	if err := self.prestart(); err == nil {
		var out, err = self.Cmd.Output()
		self.exitError = err
		self.updateStatus()
		return out, err
	} else {
		return nil, err
	}
}

func (self *Cmd) SetEnv(key string, value any) {
	var kv = fmt.Sprintf("%v=%s", key, typeutil.String(value))

	for i, pair := range self.Cmd.Env {
		k, _ := stringutil.SplitPair(pair, `=`)

		if k == key {
			self.Cmd.Env[i] = kv
			return
		}
	}

	self.Cmd.Env = append(self.Cmd.Env, kv)
}

func (self *Cmd) Run() error {
	defer self.killAndWait()

	if err := self.prestart(); err == nil {
		var err = self.Cmd.Run()
		self.exitError = err
		self.updateStatus()
		return err
	} else {
		return err
	}
}

func (self *Cmd) Start() error {
	if err := self.prestart(); err == nil {
		if err := self.Cmd.Start(); err != nil {
			self.exitError = err
			self.updateStatus()
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}

// Return the current status of the process.
func (self *Cmd) Status() Status {
	self.statusLock.Lock()
	defer self.statusLock.Unlock()

	return self.status
}

// Kill the running command.
func (self *Cmd) Kill() error {
	if p := self.Process; p != nil {
		return p.Kill()
	} else {
		return nil
	}
}

// Wait for the process to complete, then return the last status.
// Process must have been started using the Start() function.
func (self *Cmd) WaitStatus() Status {
	if !self.hasWaited {
		self.exitError = self.Wait()
		self.hasWaited = true
	}

	if xe := self.exitError; xe != nil {
		switch xe.Error() {
		case `exec: Wait was already called`:
			self.exitError = nil
		}
	}

	self.updateStatus()
	return self.status
}

// Implements io.Reader, sourcing data from the command's standard output.  If the command is not
// already running, it will be started.
func (self *Cmd) Read(p []byte) (int, error) {
	if self.reader == nil {
		if err := self.interceptStdout(); err != nil {
			return 0, err
		}

		if err := self.Start(); err != nil {
			return 0, err
		}
	}

	return self.reader.Read(p)
}

// Implements io.Writer, writing data to the commands standard input.  If the command is not already
// running, it will be started.
func (self *Cmd) Write(p []byte) (int, error) {
	if self.writer == nil {
		if err := self.interceptStdout(); err != nil {
			return 0, err
		}

		if err := self.interceptStdin(); err != nil {
			return 0, err
		}

		if err := self.Start(); err != nil {
			return 0, err
		}
	}

	return self.writer.Write(p)
}

// Implements io.Closer, killing the underlying process, waiting for it to exit, then returning.
func (self *Cmd) Close() error {
	var merr error

	if w := self.writer; w != nil {
		merr = log.AppendError(merr, w.Close())
	}

	if r := self.reader; r != nil {
		merr = log.AppendError(merr, r.Close())
	}

	merr = log.AppendError(merr, self.killAndWait())
	return merr
}

// Notify the command that no further standard input will be written.
func (self *Cmd) CloseInput() error {
	if w := self.writer; w != nil {
		return w.Close()
	} else {
		return nil
	}
}

func (self *Cmd) interceptStdout() (err error) {
	self.reader, err = self.StdoutPipe()
	return
}

func (self *Cmd) interceptStdin() (err error) {
	self.writer, err = self.StdinPipe()
	return
}

func (self *Cmd) killAndWait() error {
	// this quits the startMonitoringCommand() monitoring loop
	select {
	case self.finished <- true:
		// this waits for startMonitoringCommand() to fire callbacks and actually exit
		self.reallyDone.Wait()
	default:
	}

	if status := self.finalStatus; status != nil {
		switch msg := typeutil.String(status.Error); msg {
		case `signal: killed`:
			return nil
		default:
			return status.Error
		}
	} else {
		return nil
	}
}

// a goroutine that is launched whenever a command is started
func (self *Cmd) startMonitoringCommand() {
	self.reallyDone.Add(1)

	// if there's a timeout set, and its shorter than our monitor interval,
	// reduce the monitor interval to that
	if self.Timeout > 0 && self.Timeout < self.MonitorInterval {
		self.MonitorInterval = self.Timeout
	}

	var ticker = time.NewTicker(self.MonitorInterval)

MonitorLoop:
	for self.Timeout == 0 || time.Since(self.status.StartedAt) < self.Timeout {
		select {
		case <-ticker.C:
			self.updateStatus()

			if !self.status.Running {
				break MonitorLoop
			}

			if fn := self.OnMonitor; fn != nil {
				fn(self.Status())
			}
		case <-self.finished:
			self.Kill()
			break MonitorLoop
		}
	}

	ticker.Stop()
	self.updateStatus()
	var final = self.WaitStatus()
	self.finalStatus = &final

	// fire off callbacks
	if fn := self.OnComplete; fn != nil {
		fn(final)
	}

	if self.finalStatus.Successful {
		if fn := self.OnSuccess; fn != nil {
			fn(final)
		}
	} else if fn := self.OnError; fn != nil {
		fn(final)
	}

	self.reallyDone.Done()
}

func (self *Cmd) updateStatus() {
	self.statusLock.Lock()
	defer self.statusLock.Unlock()
	self.status.Cmd = self

	var pstate *os.ProcessState

	// handle known exit errors
	if self.exitError != nil {
		self.status.Error = self.exitError
		self.status.Running = false
		self.status.StoppedAt = time.Now()

		if xerr, ok := self.exitError.(*exec.ExitError); ok {
			pstate = xerr.ProcessState
		} else {
			return
		}
	} else if s := self.Cmd.ProcessState; s != nil {
		pstate = s
	}

	// check ProcessState
	if pstate != nil {
		self.status.PID = pstate.Pid()
		self.status.Successful = pstate.Success()

		if status, ok := pstate.Sys().(syscall.WaitStatus); ok {
			if status.Exited() {
				self.status.Running = false
				self.status.ExitCode = status.ExitStatus()
				self.status.StoppedAt = time.Now()

				if status.Signaled() {
					if sig := status.StopSignal(); sig != 0 {
						self.status.Error = errors.New(sig.String())
					} else if sig := status.Signal(); sig != 0 {
						self.status.Error = errors.New(sig.String())
					}
				} else if e := self.status.ExitCode; e != 0 {
					self.status.Error = fmt.Errorf("Process exited with status %d", e)
				}
			} else if status.Stopped() {
				self.status.Running = false
			} else {
				self.status.Running = true
			}
		}
	} else if process := self.Cmd.Process; process != nil {
		self.status.PID = process.Pid
	}
}

func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	return &Cmd{
		Cmd: exec.CommandContext(ctx, name, arg...),
	}
}
