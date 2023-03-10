package run

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
)

// LoggingRunner runs commands in the current working directory
// and streams the command output to the CLI.
// It is used by Git Town commands to run Git commands that show up in their output.
type LoggingRunner struct {
	dryRun *DryRun
	git    git
}

// NewLoggingRunner provides LoggingRunner instances.
func NewLoggingRunner(git git, dryRun *DryRun) *LoggingRunner {
	return &LoggingRunner{dryRun: dryRun, git: git}
}

// WorkingDir provides the directory that this Runner operates in.
func (r LoggingRunner) WorkingDir() string {
	return "."
}

// Run runs the given command in this ShellRunner's directory.
func (r LoggingRunner) Run(cmd string, args ...string) (*Result, error) {
	err := r.PrintCommand(cmd, args...)
	if err != nil {
		return nil, err
	}
	if r.dryRun.IsActive() {
		if len(args) == 2 && cmd == "git" && args[0] == "checkout" {
			r.dryRun.ChangeBranch(args[1])
		}
		return nil, nil //nolint:nilnil  // Can return nil result if dryRun is enabled
	}
	// Windows commands run inside CMD
	// because opening browsers is done via "start"
	if runtime.GOOS == "windows" {
		args = append([]string{"/C", cmd}, args...)
		cmd = "cmd"
	}
	subProcess := exec.Command(cmd, args...) // #nosec
	subProcess.Stderr = os.Stderr
	subProcess.Stdin = os.Stdin
	subProcess.Stdout = os.Stdout
	return nil, subProcess.Run()
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r LoggingRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		_, err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (r LoggingRunner) RunString(fullCmd string) (*Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}

// PrintCommand prints the given command-line operation on the console.
func (r LoggingRunner) PrintCommand(cmd string, args ...string) error {
	header := cmd + " "
	for index, part := range args {
		if strings.Contains(part, " ") {
			part = `"` + part + `"`
		}
		if index != 0 {
			header += " "
		}
		header += part
	}
	if cmd == "git" && r.git.IsRepository() {
		currentBranch, err := r.git.CurrentBranch()
		if err != nil {
			return err
		}
		header = fmt.Sprintf("[%s] %s", currentBranch, header)
	}
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	if err != nil {
		fmt.Println(header)
	}
	return nil
}

// PrintCommand prints the given command-line operation on the console.
func (r LoggingRunner) PrintCommandAndOutput(result *Result) error {
	err := r.PrintCommand(result.Command, result.Args...)
	fmt.Println(result.Output)
	return err
}

// git defines the Git commands needed by the LoggingShell.
type git interface {
	IsRepository() bool
	CurrentBranch() (string, error)
}
