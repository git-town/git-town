package script

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Originate/git-town/lib/git"

	"github.com/fatih/color"
	hub "github.com/github/hub/utils"
)

// OpenBrowser opens the default browser with the given URL.
func OpenBrowser(url string) {
	command, err := hub.BrowserLauncher()
	if err != nil {
		log.Fatal(err)
	}
	// command := util.GetOpenBrowserCommand()
	command = append(command, url)
	err = RunCommand(command...)

	if err != nil {
		log.Fatal(err)
	}
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(cmd ...string) {
	header := ""
	for index, part := range cmd {
		if strings.Contains(part, " ") {
			part = "\"" + strings.Replace(part, "\"", "\\\"", -1) + "\""
		}
		if index != 0 {
			header = header + " "
		}
		header = header + part
	}
	if strings.HasPrefix(header, "git") {
		header = fmt.Sprintf("[%s] %s", git.GetCurrentBranchName(), header)
	}
	fmt.Println()
	color.New(color.Bold).Println(header)
}

// RunCommand executes the given command-line operation.
func RunCommand(cmd ...string) error {
	PrintCommand(cmd...)
	subProcess := exec.Command(cmd[0], cmd[1:]...)
	subProcess.Stderr = os.Stderr
	subProcess.Stdin = os.Stdin
	subProcess.Stdout = os.Stdout
	return subProcess.Run()
}

// RunCommandSafe executes the given command-line operation, exiting if the command errors
func RunCommandSafe(cmd ...string) {
	err := RunCommand(cmd...)
	if err != nil {
		log.Fatal(err)
	}
}
