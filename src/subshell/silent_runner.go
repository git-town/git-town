package subshell

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
)

// SilentRunner runs commands in the current working directory.
// Unlike LoggingRunner, SilentRunner does not print anything to the CLI.
type SilentRunner struct {
	// When enabled, outputs the shell commands that Git Town normally runs silently
	// to the CLI.
	Debug *bool
}

// WorkingDir provides the directory that this Shell operates in.
func (r SilentRunner) WorkingDir() string {
	return "."
}

func (r SilentRunner) PrintHeader(cmd string, args ...string) {
	text := "(debug) " + cmd + " " + strings.Join(args, " ")
	_, err := color.New(color.Bold).Println(text)
	if err != nil {
		fmt.Println(text)
	}
}

func (r SilentRunner) PrintResult(text string) {
	fmt.Println(text)
}

// Run runs the given command in this ShellRunner's directory.
func (r SilentRunner) Run(cmd string, args ...string) (*Result, error) {
	if *r.Debug {
		r.PrintHeader(cmd, args...)
	}
	result, err := Exec(cmd, args...)
	if *r.Debug && result != nil {
		r.PrintResult(result.Output)
	}
	return result, err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r SilentRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (r SilentRunner) RunString(fullCmd string) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}
