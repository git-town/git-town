package test

import (
	"fmt"
	"os/exec"

	"github.com/git-town/git-town/v7/src/subshell"
)

// InternalDirRunner is an InternalRunner implementation that executes in the given directory.
type InternalDirRunner struct {
	Dir string
}

func (r InternalDirRunner) Run(executable string, args ...string) (*subshell.Output, error) {
	subProcess := exec.Command(executable, args...) // #nosec
	subProcess.Dir = r.Dir
	output, err := subProcess.CombinedOutput()
	if err != nil {
		err = subshell.ErrorDetails(executable, args, err, output)
	}
	return subshell.NewOutput(output), err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r InternalDirRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}
