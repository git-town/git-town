// Package browser allows interacting with the default browser on the user's machine.
package browser

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/git-town/git-town/v7/src/subshell"
)

// OpenBrowserCommand provides the console command to open the default browser.
func OpenBrowserCommand() string {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return "start"
	}
	openBrowserCommands := []string{
		"wsl-open",           // for Windows Subsystem for Linux, see https://github.com/git-town/git-town/issues/1344
		"garcon-url-handler", // opens links in the native browser from Crostini on ChromeOS
		"xdg-open",
		"open",
		"cygstart",
		"x-www-browser",
		"firefox",
		"opera",
		"mozilla",
		"netscape",
	}
	for _, browserCommand := range openBrowserCommands {
		fullpath, err := exec.LookPath(browserCommand)
		if err == nil && fullpath != "" {
			return fullpath
		}
	}
	return ""
}

// Open opens a new window/tab in the default browser with the given URL.
// If no browser is found, it prints the URL.
func Open(url string, runner subshell.Runner) {
	command := OpenBrowserCommand()
	if command == "" {
		fmt.Println("Please open in a browser: " + url)
		return
	}
	_, err := runner.Run(command, url)
	if err != nil {
		fmt.Println("Please open in a browser: " + url)
	}
}
