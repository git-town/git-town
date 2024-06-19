package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// FrontendRunner executes frontend shell commands.
type FrontendRunner struct {
	Backend          gitdomain.Querier
	CommandsCounter  Mutable[gohacks.Counter]
	GetCurrentBranch GetCurrentBranchFunc
	OmitBranchNames  bool
	PrintCommands    bool
}

type GetCurrentBranchFunc func(gitdomain.Querier) (gitdomain.LocalBranchName, error)

func FormatCommand(currentBranch gitdomain.LocalBranchName, omitBranch bool, executable string, args ...string) string {
	var result string
	if executable == "git" && !omitBranch {
		result = "[" + currentBranch.String() + "] git "
	} else {
		result = executable + " "
	}
	for index, part := range args {
		if part == "" {
			part = `""`
		} else if strings.Contains(part, " ") {
			part = `"` + part + `"`
		}
		if index != 0 {
			result += " "
		}
		result += part
	}
	return result
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(branch gitdomain.LocalBranchName, omitBranch bool, cmd string, args ...string) {
	header := FormatCommand(branch, omitBranch, cmd, args...)
	fmt.Println()
	fmt.Println(colors.Bold().Styled(header))
}

// Run runs the given command in this ShellRunner's directory.
func (self *FrontendRunner) Run(cmd string, args ...string) (err error) {
	*self.CommandsCounter.Value++
	var branchName gitdomain.LocalBranchName
	if !self.OmitBranchNames {
		branchName, err = self.GetCurrentBranch(self.Backend)
		if err != nil {
			return err
		}
	}
	if self.PrintCommands {
		PrintCommand(branchName, self.OmitBranchNames, cmd, args...)
	}
	if runtime.GOOS == "windows" && cmd == "start" {
		args = append([]string{"/C", cmd}, args...)
		cmd = "cmd"
	}
	subProcess := exec.Command(cmd, args...) // #nosec
	subProcess.Stderr = os.Stderr
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
	return subProcess.Wait()
}
