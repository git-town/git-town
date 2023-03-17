package subshell

import (
	"fmt"
	"os/exec"

	"github.com/kballard/go-shellquote"
)

// InternalRunner runs internal shell commands in the given working directory.
type InternalRunner struct {
	WorkingDir string
}

func (r InternalRunner) Dir() string {
	return r.WorkingDir
}

func (r InternalRunner) Run(executable string, args ...string) (*Output, error) {
	subProcess := exec.Command(executable, args...) // #nosec
	subProcess.Dir = r.WorkingDir
	output, err := subProcess.CombinedOutput()
	return NewOutput(output), err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r InternalRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

func (r InternalRunner) RunString(fullCmd string) (*Output, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}
