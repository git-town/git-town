package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/git-town/git-town/src/cli"
)

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {

	// Dir contains the directory in which to execute the command.
	// If empty, runs in the current directory.
	Dir string

	// Env allows to override the environment variables to use in the subshell, in the format provided by os.Environ()
	// If empty, uses the environment variables of this process.
	Env []string

	// Essential indicates whether this is an essential command.
	// Essential commands are critically important for Git Town to function. If they fail Git Town ends right there.
	Essential bool

	// Input contains the user input to enter into the running command.
	// It is written to the subprocess one element at a time, with a delay defined by command.InputDelay in between.
	Input []string // input into the subprocess
}

// InputDelay defines how long to wait before writing the next input string into the subprocess.
const InputDelay = 500 * time.Millisecond

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
	return &result, err
}
