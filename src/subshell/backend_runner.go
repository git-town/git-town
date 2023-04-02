package subshell

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/acarl005/stripansi"
)

// BackendRunner executes backend shell commands without output to the CLI.
type BackendRunner struct {
	// If set, runs the commands in the given directory.
	// If not set, runs the commands in the current working directory.
	Dir        *string
	Statistics Statistics
}

func (r BackendRunner) Run(executable string, args ...string) (string, error) {
	subProcess := exec.Command(executable, args...) // #nosec
	if r.Dir != nil {
		subProcess.Dir = *r.Dir
	}
	output, err := subProcess.CombinedOutput()
	if err != nil {
		err = ErrorDetails(executable, args, err, output)
	}
	return strings.TrimSpace(stripansi.Strip(string(output))), err
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
