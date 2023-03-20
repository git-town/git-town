package subshell

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
)

// InternalDebuggingRunner runs internal shell commands in the given working directory.
// It logs the executed commands and their output on the CLI.
type InternalDebuggingRunner struct {
	InternalRunner
}

func (r InternalDebuggingRunner) PrintHeader(cmd string, args ...string) {
	text := "\n(debug) " + cmd + " " + strings.Join(args, " ")
	_, err := color.New(color.Bold).Println(text)
	if err != nil {
		fmt.Println(text)
	}
}

// Run runs the given command in this ShellRunner's directory.
func (r InternalDebuggingRunner) Run(cmd string, args ...string) (*Output, error) {
	r.PrintHeader(cmd, args...)
	output, err := r.InternalRunner.Run(cmd, args...)
	if output != nil {
		fmt.Println(output)
	}
	return output, err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r InternalDebuggingRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (r InternalDebuggingRunner) RunString(fullCmd string) (*Output, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}
