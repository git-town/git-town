package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
)

// BackendRunner executes backend shell commands without output to the CLI.
type BackendRunner struct {
	CommandsCounter *gohacks.Counter
	// If set, runs the commands in the given directory.
	// If not set, runs the commands in the current working directory.
	Dir *string
	// whether to print the executed commands to the CLI
	Verbose bool
}

func (self BackendRunner) Query(executable string, args ...string) (string, error) {
	output, err := self.execute(executable, args...)
	return string(output), err
}

func (self BackendRunner) QueryTrim(executable string, args ...string) (string, error) {
	output, err := self.execute(executable, args...)
	return strings.TrimSpace(stripansi.Strip(string(output))), err
}

func (self BackendRunner) Run(executable string, args ...string) error {
	_, err := self.execute(executable, args...)
	return err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (self BackendRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := self.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf(messages.RunCommandProblem, argv, err)
		}
	}
	return nil
}

func (self BackendRunner) execute(executable string, args ...string) ([]byte, error) {
	self.CommandsCounter.Register()
	if self.Verbose {
		printHeader(executable, args...)
	}
	subProcess := exec.Command(executable, args...) // #nosec
	if self.Dir != nil {
		subProcess.Dir = *self.Dir
	}
	subProcess.Env = append(subProcess.Environ(), "LC_ALL=C")
	outputBytes, err := subProcess.CombinedOutput()
	if err != nil {
		err = ErrorDetails(executable, args, err, outputBytes)
	}
	if self.Verbose && len(outputBytes) > 0 {
		os.Stdout.Write(outputBytes)
	}
	return outputBytes, err
}

func ErrorDetails(executable string, args []string, err error, output []byte) error {
	return fmt.Errorf(`
----------------------------------------
Diagnostic information of failed command

COMMAND: %s %v
ERROR: %w
OUTPUT START
%s
OUTPUT END
----------------------------------------`, executable, strings.Join(args, " "), err, string(output))
}

func printHeader(cmd string, args ...string) {
	quoted := stringslice.SurroundEmptyWith(args, `"`)
	text := "\n(verbose) " + cmd + " " + strings.Join(quoted, " ")
	fmt.Println(colors.Bold().Styled(text))
}
