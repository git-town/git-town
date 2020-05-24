package command

import (
	"fmt"
	"os"

	"github.com/kballard/go-shellquote"
)

// ShellInCurrentDir is an implementation of the Shell interface that runs commands in the current working directory.
type ShellInCurrentDir struct {
	dir string
}

// NewShellInCurrentDir provides ShellInCurrentDir instances.
func NewShellInCurrentDir() (ShellInCurrentDir, error) {
	dir, err := os.Getwd()
	return ShellInCurrentDir{dir}, err
}

// Dir provides the directory that this Shell operates in.
func (shell ShellInCurrentDir) Dir() string {
	return shell.dir
}

// MustRun runs the given command and returns the result. Panics on error.
func (shell ShellInCurrentDir) MustRun(cmd string, args ...string) *Result {
	return MustRun(cmd, args...)
}

// Run runs the given command in this ShellRunner's directory.
func (shell ShellInCurrentDir) Run(cmd string, args ...string) (*Result, error) {
	return Run(cmd, args...)
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (shell ShellInCurrentDir) RunMany(commands [][]string) error {
	for _, argv := range commands {
		outcome, err := Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w\n%v", argv, err, outcome)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell ShellInCurrentDir) RunString(fullCmd string) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return Run(cmd, args...)
}

// RunStringWith runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell ShellInCurrentDir) RunStringWith(fullCmd string, options Options) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return RunWith(options, cmd, args...)
}
