package script

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/util"

	"github.com/fatih/color"
)

var browserTools = []string{"xdg-open", "open"}
var missingBrowserToolMessages = []string{
	"Opening a browser requires 'open' on Mac or 'xdg-open' on Linux.",
	"If you would like another command to be supported,",
	"please open an issue at https://github.com/Originate/git-town/issues",
}

func OpenBrowser(url string) {
	for _, browserTool := range browserTools {
		if util.GetCommandOutput("which", browserTool) != "" {
			err := RunCommand(browserTool, url)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	util.ExitWithErrorMessage(missingBrowserToolMessages...)
}

func PrintCommand(cmd ...string) {
	header := strings.Join(cmd, " ")
	if strings.HasPrefix(header, "git") {
		header = fmt.Sprintf("[%s] %s", git.GetCurrentBranchName(), header)
	}
	fmt.Println()
	color.New(color.Bold).Println(header)
}

func RunCommand(cmd ...string) error {
	PrintCommand(cmd...)
	subProcess := exec.Command(cmd[0], cmd[1:]...)
	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stderr
	return subProcess.Run()
}
