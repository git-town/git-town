package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/cache"
)

// FrontendRunner executes frontend shell commands.
type FrontendRunner struct {
	CurrentBranch   *cache.String
	OmitBranchNames bool
	Stats           Statistics
}

// Run runs the given command in this ShellRunner's directory.
func (r *FrontendRunner) Run(cmd string, args ...string) error {
	r.Stats.RegisterRun()
	var branchName string
	if r.OmitBranchNames {
		branchName = ""
	} else {
		branchName = r.CurrentBranch.Value()
	}
	PrintCommand(branchName, r.OmitBranchNames, cmd, args...)
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
func (r *FrontendRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(branch string, omitBranch bool, cmd string, args ...string) {
	header := FormatCommand(branch, omitBranch, cmd, args...)
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	if err != nil {
		fmt.Println(header)
	}
}

func FormatCommand(currentBranch string, omitBranch bool, executable string, args ...string) string {
	result := ""
	if executable == "git" && !omitBranch {
		result = "[" + currentBranch + "] git "
	} else {
		result = executable + " "
	}
	for index, part := range args {
		if strings.Contains(part, " ") {
			part = `"` + part + `"`
		}
		if index != 0 {
			result += " "
		}
		result += part
	}
	return result
}
