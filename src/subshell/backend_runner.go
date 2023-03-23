package subshell

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
)

// BackendRunner executes backend shell commands.
type BackendRunner struct {
	// If set, runs the commands in the given directory.
	// If not set, runs the commands in the current working directory.
	Dir *string
}

func (r BackendRunner) Run(executable string, args ...string) (*Output, error) {
	subProcess := exec.Command(executable, args...) // #nosec
	if r.Dir != nil {
		subProcess.Dir = *r.Dir
	}
	output, err := subProcess.CombinedOutput()
	if err != nil {
		err = ErrorDetails(executable, args, err, output)
	}
	return NewOutput(output), err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r BackendRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

func (r BackendRunner) RunString(fullCmd string) (*Output, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}

func ErrorDetails(executable string, args []string, err error, output []byte) error {
	return fmt.Errorf(`
----------------------------------------
Diagnostic information of failed command

Command: %s %v
Error: %w
Output:
%s
----------------------------------------`, executable, strings.Join(args, " "), err, string(output))
}
