package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/messages"
)

// BackendRunner executes backend shell commands without output to the CLI.
type BackendRunner struct {
	// If set, runs the commands in the given directory.
	// If not set, runs the commands in the current working directory.
	Dir   *string
	Stats Statistics
	// whether to print the executed commands to the CLI
	Verbose bool
}

func (r BackendRunner) Query(executable string, args ...string) (string, error) {
	output, err := r.execute(executable, args...)
	return string(output), err
}

func (r BackendRunner) QueryTrim(executable string, args ...string) (string, error) {
	output, err := r.execute(executable, args...)
	return strings.TrimSpace(stripansi.Strip(string(output))), err
}

func (r BackendRunner) Run(executable string, args ...string) error {
	_, err := r.execute(executable, args...)
	return err
}

func (r BackendRunner) execute(executable string, args ...string) ([]byte, error) {
	r.Stats.RegisterRun()
	if r.Verbose {
		printHeader(executable, args...)
	}
	subProcess := exec.Command(executable, args...) // #nosec
	if r.Dir != nil {
		subProcess.Dir = *r.Dir
	}
	outputBytes, err := subProcess.CombinedOutput()
	if err != nil {
		err = ErrorDetails(executable, args, err, outputBytes)
	}
	if r.Verbose && len(outputBytes) > 0 {
		os.Stdout.Write(outputBytes)
	}
	return outputBytes, err
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r BackendRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf(messages.RunProblem, argv, err)
		}
	}
	return nil
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
	text := "\n(debug) " + cmd + " " + strings.Join(args, " ")
	_, err := color.New(color.Bold).Println(text)
	if err != nil {
		fmt.Println(text)
	}
}
