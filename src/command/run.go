package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {
	Dir   string   // the directory in which to execute the command
	Env   []string // environment variables to use, in the format provided by os.Environ()
	Input []string // input into the subprocess
}

// MustRun executes an essential subshell command given in argv notation.
// Essential subshell commands are essential for the functioning of Git Town.
// If they fail, Git Town ends right there.
func MustRun(cmd string, args ...string) *Result {
	result, err := RunWith(Options{}, cmd, args...)
	if err != nil {
		fmt.Printf("\n\nError running '%s %s': %s", cmd, strings.Join(args, " "), err)
		os.Exit(1)
	}
	return result
}

// MustRunInDir executes an essential subshell command given in argv notation.
// Essential subshell commands are essential for the functioning of Git Town.
// If they fail, Git Town ends right there.
func MustRunInDir(dir string, cmd string, args ...string) *Result {
	result, err := RunWith(Options{Dir: dir}, cmd, args...)
	if err != nil {
		fmt.Printf("\n\nError running '%s %s' in %s: %s", cmd, strings.Join(args, " "), dir, err)
		os.Exit(1)
	}
	return result
}

// MustRunWith runs an essential subshell command with the given options.
// Essential subshell commands are essential for the functioning of Git Town.
// If they fail, Git Town ends right there.
func MustRunWith(opts Options, cmd string, args ...string) *Result {
	result, err := RunWith(opts, cmd, args...)
	if err != nil {
		fmt.Printf("\n\nError running with options %v: %v", opts, err)
		os.Exit(1)
	}
	return result
}

// Run executes the command given in argv notation.
func Run(cmd string, args ...string) (*Result, error) {
	return RunWith(Options{}, cmd, args...)
}

// RunInDir executes the given command in the given directory.
func RunInDir(dir string, cmd string, args ...string) (*Result, error) {
	return RunWith(Options{Dir: dir}, cmd, args...)
}

// RunWith runs the command with the given RunOptions.
func RunWith(opts Options, cmd string, args ...string) (*Result, error) {
	logRun(cmd, args...)
	subProcess := exec.Command(cmd, args...) // #nosec
	if opts.Dir != "" {
		subProcess.Dir = opts.Dir
	}
	if opts.Env != nil {
		subProcess.Env = opts.Env
	}
	var output bytes.Buffer
	subProcess.Stdout = &output
	subProcess.Stderr = &output
	result := Result{command: cmd, args: args}
	input, err := subProcess.StdinPipe()
	if err != nil {
		return &result, err
	}
	err = subProcess.Start()
	if err != nil {
		return &result, fmt.Errorf("can't start subprocess '%s %s': %w", cmd, strings.Join(args, " "), err)
	}
	for _, userInput := range opts.Input {
		// Here we simply wait for some time until the subProcess needs the input.
		// Capturing the output and scanning for the actual content needed
		// would introduce substantial amounts of multi-threaded complexity
		// for not enough gains.
		time.Sleep(50 * time.Millisecond)
		_, err := input.Write([]byte(userInput))
		if err != nil {
			result.output = output.String()
			return &result, errors.Wrapf(err, "can't write %q to subprocess '%s %s'", userInput, cmd, strings.Join(args, " "))
		}
	}
	err = subProcess.Wait()
	result.output = output.String()
	return &result, err
}
