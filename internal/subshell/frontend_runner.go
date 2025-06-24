package subshell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// FrontendRunner executes frontend shell commands.
type FrontendRunner struct {
	Backend          subshelldomain.Querier
	CommandsCounter  Mutable[gohacks.Counter]
	GetCurrentBranch GetCurrentBranchFunc
	PrintBranchNames bool
	PrintCommands    bool
}

type GetCurrentBranchFunc func(subshelldomain.Querier) (gitdomain.LocalBranchName, error)

func FormatCommand(currentBranch gitdomain.LocalBranchName, printBranch bool, env []string, executable string, args ...string) string {
	result := ""
	if printBranch {
		result += "[" + currentBranch.String() + "] "
	}
	if len(env) > 0 {
		result += strings.Join(env, " ") + " "
	}
	result += executable + " "
	quoted := stringslice.SurroundEmptyWith(args, `"`)
	quoted = stringslice.SurroundSpacesWith(quoted, `"`)
	result += strings.Join(quoted, " ")
	return result
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(branch gitdomain.LocalBranchName, printBranch bool, env []string, cmd string, args ...string) {
	header := FormatCommand(branch, printBranch, env, cmd, args...)
	fmt.Println()
	fmt.Println(colors.Bold().Styled(header))
}

func (self *FrontendRunner) Run(cmd string, args ...string) error {
	return self.execute([]string{}, cmd, args...)
}

func (self *FrontendRunner) RunWithEnv(env []string, cmd string, args ...string) error {
	return self.execute(env, cmd, args...)
}

// runs the given command in this ShellRunner's directory.
func (self *FrontendRunner) execute(env []string, cmd string, args ...string) (err error) {
	self.CommandsCounter.Value.Inc()
	var branchName gitdomain.LocalBranchName
	if self.PrintBranchNames {
		branchName, err = self.GetCurrentBranch(self.Backend)
		if err != nil {
			return err
		}
	}
	if self.PrintCommands {
		PrintCommand(branchName, self.PrintBranchNames, env, cmd, args...)
	}
	if runtime.GOOS == "windows" && cmd == "start" {
		args = append([]string{"/C", cmd}, args...)
		cmd = "cmd"
	}
	concurrentGitRetriesLeft := concurrentGitRetries
	for {
		subProcess := exec.Command(cmd, args...)
		subProcess.Env = append(subProcess.Environ(), env...)
		var stderrBuffer bytes.Buffer // we only need to look at STDERR since that's where Git will print error messages
		subProcess.Stderr = io.MultiWriter(os.Stderr, &stderrBuffer)
		subProcess.Stdin = os.Stdin
		subProcess.Stdout = os.Stdout
		err = subProcess.Start()
		if err != nil {
			return err
		}
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT) // Listen for Ctrl-C
		go func() {
			<-interrupt
			if err := subProcess.Process.Release(); err == nil {
				// process has already finished, no need to kill
				return
			}
			fmt.Printf("Abort detected, shutting down %q gracefully ...", strings.Join(append([]string{cmd}, args...), " "))
			if err := subProcess.Process.Kill(); err != nil {
				fmt.Println("Error killing subprocess:", err)
			}
		}()
		err = subProcess.Wait()
		if err == nil {
			break
		}
		if !containsConcurrentGitAccess(stderrBuffer.String()) {
			break
		}
		concurrentGitRetriesLeft -= 1
		if concurrentGitRetriesLeft == 0 {
			break
		}
		fmt.Println(colors.Bold().Styled("\n" + messages.GitAnotherProcessIsRunningRetry + "\n"))
		time.Sleep(concurrentGitRetryDelay)
	}
	return err
}
