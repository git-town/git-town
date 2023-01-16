package run

import (
	"fmt"

	"github.com/kballard/go-shellquote"
)

// SilentShell is an implementation of the Shell interface that runs commands in the current working directory.
type SilentShell struct{}

// WorkingDir provides the directory that this Shell operates in.
func (s SilentShell) WorkingDir() string {
	return "."
}

// Run runs the given command in this ShellRunner's directory.
func (s SilentShell) Run(cmd string, args ...string) (*Result, error) {
	return Exec(cmd, args...)
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (s SilentShell) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := s.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (s SilentShell) RunString(fullCmd string) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return s.Run(cmd, args...)
}

// RunStringWith runs the given command (including possible arguments) in this ShellInDir's directory.
func (s SilentShell) RunStringWith(fullCmd string, options *Options) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return WithOptions(options, cmd, args...)
}
