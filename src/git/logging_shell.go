package git

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/dryrun"
	"github.com/kballard/go-shellquote"
)

// LoggingShell is an implementation of the Shell interface
// that runs commands in the current working directory
// and streams the command output to the application output.
// It is used by Git Town commands to run Git commands that show up in their output.
type LoggingShell struct {
	silentRunner *Runner
}

// NewLoggingShell provides StreamingShell instances.
func NewLoggingShell(silent *Runner) *LoggingShell {
	return &LoggingShell{silent}
}

// WorkingDir provides the directory that this Shell operates in.
func (shell LoggingShell) WorkingDir() string {
	return "."
}

// Run runs the given command in this ShellRunner's directory.
func (shell LoggingShell) Run(cmd string, args ...string) (*command.Result, error) {
	err := shell.PrintCommand(cmd, args...)
	if err != nil {
		return nil, err
	}
	if dryrun.IsActive() {
		if len(args) == 2 && cmd == "git" && args[0] == "checkout" {
			dryrun.SetCurrentBranchName(args[1])
		}
		return nil, nil
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
func (shell LoggingShell) RunMany(commands [][]string) error {
	for _, argv := range commands {
		outcome, err := shell.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w\n%v", argv, err, outcome)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell LoggingShell) RunString(fullCmd string) (*command.Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return shell.Run(cmd, args...)
}

// RunStringWith runs the given command (including possible arguments) in this ShellInDir's directory.
func (shell LoggingShell) RunStringWith(fullCmd string, options command.Options) (*command.Result, error) {
	panic("this isn't used")
}

// PrintCommand prints the given command-line operation on the console.
func (shell LoggingShell) PrintCommand(cmd string, args ...string) error {
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
	if cmd == "git" && shell.silentRunner.IsRepository() {
		currentBranch, err := shell.silentRunner.CurrentBranch()
		if err != nil {
			return err
		}
		header = fmt.Sprintf("[%s] %s", currentBranch, header)
	}
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	if err != nil {
		panic(err)
	}
	return nil
}
