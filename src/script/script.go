package script

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/git-town/git-town/src/browsers"
	"github.com/git-town/git-town/src/dryrun"
	"github.com/git-town/git-town/src/git"

	"github.com/fatih/color"
)

var dryRunMessage = `
In dry run mode. No commands will be run. When run in normal mode, the command
output will appear beneath the command. Some commands will only be run if
necessary. For example: 'git push' will run if and only if there are local
commits not on the remote.
`

// ActivateDryRun causes all commands to not be run
func ActivateDryRun() {
	_, err := color.New(color.FgBlue).Print(dryRunMessage)
	if err != nil {
		panic(err)
	}
	dryrun.Activate(git.GetCurrentBranchName())
}

// OpenBrowser opens the default browser with the given URL.
// If no browser is found, prints the URL.
func OpenBrowser(url string) {
	command := browsers.GetOpenBrowserCommand()
	if command == "" {
		fmt.Println("Please open in a browser: " + url)
		return
	}
	err := RunCommand(command, url)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Please open in a browser: " + url)
	}
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(cmd string, args ...string) {
	header := cmd + " "
	for index, part := range args {
		if strings.Contains(part, " ") {
			part = "\"" + strings.Replace(part, "\"", "\\\"", -1) + "\""
		}
		if index != 0 {
			header += " "
		}
		header += part
	}
	if cmd == "git" && git.IsRepository() {
		header = fmt.Sprintf("[%s] %s", git.GetCurrentBranchName(), header)
	}
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	if err != nil {
		panic(err)
	}
}

// RunCommand executes the given command-line operation.
func RunCommand(cmd string, args ...string) error {
	PrintCommand(cmd, args...)
	if dryrun.IsActive() {
		if len(args) == 2 && cmd == "git" && args[0] == "checkout" {
			dryrun.SetCurrentBranchName(args[1])
		}
		return nil
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
	return subProcess.Run()
}

// RunCommandSafe executes the given command-line operation, exiting if the command errors
func RunCommandSafe(cmd string, args ...string) {
	err := RunCommand(cmd, args...)
	if err != nil {
		panic(err)
	}
}
