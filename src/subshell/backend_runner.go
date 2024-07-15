package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
)

// BackendRunner executes backend shell commands without output to the CLI.
type BackendRunner struct {
	CommandsCounter Mutable[gohacks.Counter]
	// If set, runs the commands in the given directory.
	// If not set, runs the commands in the current working directory.
	Dir Option[string]
	// whether to print the executed commands to the CLI
	Verbose configdomain.Verbose
}

func (self BackendRunner) Query(executable string, args ...string) (string, error) {
	return self.execute(executable, args...)
}

func (self BackendRunner) QueryTrim(executable string, args ...string) (string, error) {
	output, err := self.execute(executable, args...)
	return strings.TrimSpace(stripansi.Strip(output)), err
}

func (self BackendRunner) Run(executable string, args ...string) error {
	_, err := self.execute(executable, args...)
	return err
}

func (self BackendRunner) execute(executable string, args ...string) (string, error) {
	self.CommandsCounter.Value.Inc()
	if self.Verbose {
		printHeader(executable, args...)
	}
	subProcess := exec.Command(executable, args...) // #nosec
	if dir, has := self.Dir.Get(); has {
		subProcess.Dir = dir
	}
	subProcess.Env = append(subProcess.Environ(), "LC_ALL=C")
	outputBytes, err := subProcess.CombinedOutput()
	if err != nil {
		err = ErrorDetails(executable, args, err, outputBytes)
	}
	if self.Verbose && len(outputBytes) > 0 {
		os.Stdout.Write(outputBytes)
	}
	return string(outputBytes), err
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
