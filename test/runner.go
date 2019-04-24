package test

import (
	"os/exec"
	"strings"
)

// Runner runs command-line apps like "git-town" in a subshell
type Runner struct {
	Output string
	Err    error
}

// Run runs the given command with the given argv-like arguments
// and stores the output and error for later analysis.
func (r *Runner) Run(name string, arguments ...string) {
	cmd := exec.Command(name, arguments...)
	rawOutput, err := cmd.CombinedOutput()
	r.Output = string(rawOutput)
	r.Err = err
}

// RunString runs the given command (that can contain arguments)
// and stores the output and error for later analysis.
//
// Currently this splits the string by space,
// this only works for simple commands without quotes.
func (r *Runner) RunString(command string) {
	parts := strings.Fields(command)
	command, args := parts[0], parts[1:]
	r.Run(command, args...)
}

// OutputContains returns whether the output of the last command
// contains the given string.
func (r *Runner) OutputContains(text string) bool {
	return strings.Contains(r.Output, text)
}
