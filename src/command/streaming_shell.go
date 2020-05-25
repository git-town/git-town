package command

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/script"
	"github.com/kballard/go-shellquote"
)

// StreamingShell is an implementation of the Shell interface
// that runs commands in the current working directory
// and streams the command output to the application output.
type StreamingShell struct {
	out *os.File // where to stream command output
}

// NewStreamingShell provides StreamingShell instances.
func NewStreamingShell(out *os.File) *StreamingShell {
	return &StreamingShell{out}
}

// WorkingDir provides the directory that this Shell operates in.
func (shell StreamingShell) WorkingDir() string {
	return "."
}

// MustRun runs the given command and returns the result. Panics on error.
func (shell StreamingShell) MustRun(cmd string, args ...string) *Result {
	return MustRun(cmd, args...)
}

// Run runs the given command in this ShellRunner's directory.
func (shell StreamingShell) Run(cmd string, args ...string) (*Result, error) {
	err := script.RunCommand(cmd, args...)
	return nil, err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (shell StreamingShell) RunMany(commands [][]string) error {
	for _, argv := range commands {
		outcome, err := Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w\n%v", argv, err, outcome)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell StreamingShell) RunString(fullCmd string) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return Run(cmd, args...)
}

// RunStringWith runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell StreamingShell) RunStringWith(fullCmd string, options Options) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return RunWith(options, cmd, args...)
}
