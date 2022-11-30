package run

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/git-town/git-town/v7/src/cli"
)

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {
	// Dir contains the directory in which to execute the command.
	// If empty, runs in the current directory.
	Dir string

	// Env allows to override the environment variables to use in the subshell, in the format provided by os.Environ()
	// If empty, uses the environment variables of this process.
	Env []string

	// Input contains the user input to enter into the running command.
	// It is written to the subprocess one element at a time, with a delay defined by command.InputDelay in between.
	Input []string // input into the subprocess
}

// InputDelay defines how long to wait before writing the next input string into the subprocess.
const InputDelay = 500 * time.Millisecond

// Exec executes the command given in argv notation.
func Exec(cmd string, args ...string) (*Result, error) {
	return WithOptions(Options{}, cmd, args...)
}

// InDir executes the given command in the given directory.
func InDir(dir string, cmd string, args ...string) (*Result, error) {
	return WithOptions(Options{Dir: dir}, cmd, args...)
}

// WithOptions runs the command with the given RunOptions.
func WithOptions(opts Options, cmd string, args ...string) (*Result, error) {
	cli.LogRun(cmd, args...)
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
		// https://github.com/git-town/go-execplus could help make this more robust.
		time.Sleep(InputDelay)
		_, err := input.Write([]byte(userInput))
		if err != nil {
			result.output = output.String()
			return &result, fmt.Errorf("can't write %q to subprocess '%s %s': %w", userInput, cmd, strings.Join(args, " "), err)
		}
	}
	err = subProcess.Wait()
	result.output = output.String()
	result.exitCode = subProcess.ProcessState.ExitCode()
	if err != nil {
		err = fmt.Errorf(`
----------------------------------------
Diagnostic information of failed command

Command: %s %s
Error: %w
Output:
%s
----------------------------------------`, cmd, strings.Join(args, " "), err, result.output)
	}
	return &result, err
}
