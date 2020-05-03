package command

import (
	"fmt"

	"github.com/kballard/go-shellquote"
)

// ShellInDir is an implementation of the Shell interface that runs commands in a given directory.
type ShellInDir struct {
	Dir string // the directory path in which this instance runs shell commands
}

// MustRun runs the given command and returns the result. Panics on error.
func (shell *ShellInDir) MustRun(cmd string, args ...string) (result *Result) {
	return MustRunInDir(shell.Dir, cmd, args...)
}

// Run runs the given command in this ShellRunner's directory.
func (shell *ShellInDir) Run(cmd string, args ...string) (result *Result, err error) {
	return RunInDir(shell.Dir, cmd, args...)
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (shell *ShellInDir) RunMany(commands [][]string) error {
	for _, argv := range commands {
		outcome, err := RunInDir(shell.Dir, argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w\n%v", argv, err, outcome)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell *ShellInDir) RunString(fullCmd string) (result *Result, err error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return result, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return RunInDir(shell.Dir, cmd, args...)
}

// RunStringWith runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell *ShellInDir) RunStringWith(fullCmd string, options Options) (result *Result, err error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return result, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	options.Dir = shell.Dir
	return RunWith(options, cmd, args...)
}
