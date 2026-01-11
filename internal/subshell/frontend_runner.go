package subshell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FrontendRunner executes frontend shell commands.
type FrontendRunner struct {
	Backend          subshelldomain.Querier
	CommandsCounter  Mutable[gohacks.Counter]
	GetCurrentBranch GetCurrentBranchFunc
	GetCurrentSHA    GetCurrentSHAFunc
	PrintBranchNames bool
	PrintCommands    bool
}

type (
	GetCurrentBranchFunc func(subshelldomain.Querier) (Option[gitdomain.LocalBranchName], error)
	GetCurrentSHAFunc    func(subshelldomain.Querier) (gitdomain.SHA, error)
)

func FormatCommand(location gitdomain.Location, printBranch bool, env []string, executable string, args ...string) string {
	result := ""
	if printBranch {
		result += "[" + location.String() + "] "
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
func PrintCommand(location gitdomain.Location, printBranch bool, env []string, cmd string, args ...string) {
	header := FormatCommand(location, printBranch, env, cmd, args...)
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
func (self *FrontendRunner) execute(env []string, cmd string, args ...string) error {
	self.CommandsCounter.Value.Increment()
	var location gitdomain.Location
	if self.PrintBranchNames {
		currentBranchOpt, err := self.GetCurrentBranch(self.Backend)
		if err != nil {
			return err
		}
		if currentBranch, has := currentBranchOpt.Get(); has {
			location = currentBranch.Location()
		} else {
			currentSHA, err := self.GetCurrentSHA(self.Backend)
			if err != nil {
				return err
			}
			location = currentSHA.Truncate(7).Location()
		}
	}
	if self.PrintCommands {
		PrintCommand(location, self.PrintBranchNames, env, cmd, args...)
	}
	if runtime.GOOS == "windows" && cmd == "start" {
		args = append([]string{"/C", cmd}, args...)
		cmd = "cmd"
	}
	concurrentGitRetriesLeft := concurrentGitRetries
	var err error
	for {
		subProcess := exec.CommandContext(context.Background(), cmd, args...)
		subProcess.Env = append(subProcess.Environ(), env...)
		var stderrBuffer bytes.Buffer // we only need to look at STDERR since that's where Git will print error messages
		subProcess.Stderr = io.MultiWriter(os.Stderr, &stderrBuffer)
		subProcess.Stdin = os.Stdin
		subProcess.Stdout = os.Stdout
		if err := subProcess.Start(); err != nil {
			return err
		}
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT) // Listen for Ctrl-C
		go func() {
			<-interrupt
			if err := subProcess.Process.Release(); err == nil {
				// process has already finished, no need to stop it
				return
			}
			fmt.Printf("Abort detected, shutting down %q gracefully ...", strings.Join(append([]string{cmd}, args...), " "))
			if err := subProcess.Process.Kill(); err != nil {
				fmt.Println("Error stopping subprocess:", err)
			}
		}()
		if err = subProcess.Wait(); err == nil {
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
