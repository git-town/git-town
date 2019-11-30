package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {

	// Dir contains the directory in which to execute the command.
	// Runs in the current directory if this option is not provided.
	Dir string

	// Env allows to override the environment variables to use in the subshell, in the format provided by os.Environ()
	// Uses the environment variables of this process if this option is not provided.
	Env []string

	// Essential indicates whether this is an essential command.
	// Essential commands are critically important for Git Town to function., if they fail Git Town ends right there.
	Essential bool
}

// MustRun executes an essential subshell command given in argv notation.
// Essential subshell commands are essential for the functioning of Git Town.
// If they fail, Git Town ends right there.
func MustRun(cmd string, args ...string) *Result {
	result, _ := RunWith(Options{Essential: true}, cmd, args...)
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
	output, err := subProcess.CombinedOutput()
	if opts.Essential && err != nil {
		fmt.Printf("\n\nError running '%s %s': %s", cmd, strings.Join(args, " "), err)
		os.Exit(1)
	}
	result := Result{
		command: cmd,
		args:    args,
		output:  string(output),
	}
	return &result, err
}
