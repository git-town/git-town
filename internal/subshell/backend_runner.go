package subshell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v16/internal/cli/colors"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/gohacks"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
	concurrentGitRetriesLeft := concurrentGitRetries
	var outputText string
	var outputBytes []byte
	var err error
	for {
		outputBytes, err = subProcess.CombinedOutput()
		outputText = string(outputBytes)
		if err == nil {
			break
		}
		if !containsConcurrentGitAccess(outputText) {
			err = ErrorDetails(executable, args, err, outputBytes)
			break
		}
		concurrentGitRetriesLeft -= 1
		if concurrentGitRetriesLeft == 0 {
			break
		}
		fmt.Println(messages.GitAnotherProcessIsRunningRetry)
		time.Sleep(concurrentGitRetryDelay)
	}
	if self.Verbose && len(outputBytes) > 0 {
		os.Stdout.Write(bytes.ReplaceAll(outputBytes, []byte{0x00}, []byte{'\n', '\n'}))
	}
	return outputText, err
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

func containsConcurrentGitAccess(text string) bool {
	return strings.Contains(text, "fatal: Unable to create '") && strings.Contains(text, "index.lock': File exists.")
}

func printHeader(cmd string, args ...string) {
	quoted := stringslice.SurroundEmptyWith(args, `"`)
	text := "\n(verbose) " + cmd + " " + strings.Join(quoted, " ")
	fmt.Println(colors.Bold().Styled(text))
}
