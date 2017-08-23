package script

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Originate/git-town/src/dryrun"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"

	"github.com/fatih/color"
)

// ActivateDryRun causes all commands to not be run
func ActivateDryRun() {
	dryrun.Activate(git.GetCurrentBranchName())
}

// OpenBrowser opens the default browser with the given URL.
func OpenBrowser(url string) {
	command := util.GetOpenBrowserCommand()
	err := RunCommand(command, url)
	if err != nil {
		log.Fatal(err)
	}
}

// PrintCommand prints the given command-line operation on the console.
func PrintCommand(cmd ...string) {
	header := formatCommand(cmd...)
	if strings.HasPrefix(header, "git") && git.IsRepository() {
		header = fmt.Sprintf("[%s] %s", git.GetCurrentBranchName(), header)
	}
	fmt.Println()
	color.New(color.Bold).Println(header)
}

// RunCommand executes the given command-line operation.
func RunCommand(cmd ...string) error {
	PrintCommand(cmd...)
	if dryrun.IsActive() {
		if len(cmd) == 3 && cmd[0] == "git" && cmd[1] == "checkout" {
			dryrun.SetCurrentBranchName(cmd[2])
		}
		fmt.Printf("< output of: %s >\n", formatCommand(cmd...))
		return nil
	}
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

func formatCommand(cmd ...string) string {
	result := ""
	for index, part := range cmd {
		if strings.Contains(part, " ") {
			part = "\"" + strings.Replace(part, "\"", "\\\"", -1) + "\""
		}
		if index != 0 {
			result = result + " "
		}
		result = result + part
	}
	return result
}
