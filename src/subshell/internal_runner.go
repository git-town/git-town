package subshell

import (
	"fmt"
	"os/exec"

	"github.com/kballard/go-shellquote"
)

// InternalRunner runs internal shell commands in the given working directory.
type InternalRunner struct{}

func (r InternalRunner) Run(dir string, executable string, args ...string) (*Output, error) {
	subProcess := exec.Command(executable, args...) // #nosec
	subProcess.Dir = dir
	output, err := subProcess.CombinedOutput()
	return NewOutput(output), err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r InternalRunner) RunMany(dir string, commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(dir, argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

func (r InternalRunner) RunString(dir, fullCmd string) (*Output, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(dir, cmd, args...)
}
