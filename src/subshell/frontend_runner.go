package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
)

// FrontendRunner executes frontend shell commands.
type FrontendRunner struct {
	GetCurrentBranch GetCurrentBranchFunc
	OmitBranchNames  bool
	CommandsCounter  *gohacks.Counter
}

type GetCurrentBranchFunc func() (domain.LocalBranchName, error)

// Run runs the given command in this ShellRunner's directory.
func (fr *FrontendRunner) Run(cmd string, args ...string) (err error) {
	fr.CommandsCounter.Register()
	var branchName domain.LocalBranchName
	if !fr.OmitBranchNames {
		branchName, err = fr.GetCurrentBranch()
		if err != nil {
			return err
		}
	}
	PrintCommand(branchName, fr.OmitBranchNames, cmd, args...)
	// Windows commands run inside CMD
	// because opening browsers is done via "start"
	// TODO: do this only when actually running the "start" command
	if runtime.GOOS == "windows" {
		args = append([]string{"/C", cmd}, args...)
		cmd = "cmd"
	}
	subProcess := exec.Command(cmd, args...) // #nosec
	subProcess.Stderr = os.Stderr
	subProcess.Stdin = os.Stdin
	subProcess.Stdout = os.Stdout
	return subProcess.Run()
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (fr *FrontendRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := fr.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf(messages.RunCommandProblem, argv, err)
		}
	}
	return nil
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(branch domain.LocalBranchName, omitBranch bool, cmd string, args ...string) {
	header := FormatCommand(branch, omitBranch, cmd, args...)
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	if err != nil {
		fmt.Println(header)
	}
}

func FormatCommand(currentBranch domain.LocalBranchName, omitBranch bool, executable string, args ...string) string {
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
