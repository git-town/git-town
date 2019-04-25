package test

import (
	"os/exec"
	"strings"
)

// Runner runs shell commands.
type Runner struct {
	environments *Environments
}

// Run runs the given command with the given argv-like arguments in the current directory
// and stores the output and error for later analysis.
func (r *Runner) Run(name string, arguments ...string) (string, error) {
	cmd := exec.Command(name, arguments...)
	rawOutput, err := cmd.CombinedOutput()
	return string(rawOutput), err
}

// RunString runs the given command (that can contain arguments) in the current directory
// and stores the output and error for later analysis.
//
// Currently this splits the string by space,
// this only works for simple commands without quotes.
func (r *Runner) RunString(command string) (string, error) {
	parts := strings.Fields(command)
	command, args := parts[0], parts[1:]
	return r.Run(command, args...)
}

// RunMany runs all given commands in current directory.
// Failed commands cause abortion of the function with the received error.
func (r *Runner) RunMany(commands [][]string) error {
	for _, commandList := range commands {
		command, args := commandList[0], commandList[1:]
		_, err := r.Run(command, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
