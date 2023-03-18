package subshell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/kballard/go-shellquote"
)

// PublicRunner executes the given shell commands and streams their output to the CLI.
type PublicRunner struct {
	CurrentBranch *cache.String
}

// Run runs the given command in this ShellRunner's directory.
func (r PublicRunner) Run(cmd string, args ...string) error {
	PrintCommand(r.CurrentBranch.Value(), cmd, args...)
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
func (r PublicRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (r PublicRunner) RunString(fullCmd string) error {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(branch string, cmd string, args ...string) {
	header := FormatCommand(branch, cmd, args...)
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	if err != nil {
		fmt.Println(header)
	}
}

func FormatCommand(currentBranch string, executable string, args ...string) string {
	result := ""
	if executable == "git" {
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
