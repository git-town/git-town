package browsers

import (
	"fmt"
	"runtime"

	"github.com/git-town/git-town/src/command"
)

// OpenBrowserCommand returns the command to run on the console
// to open the default browser.
func OpenBrowserCommand() string {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters
		//       like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return "start"
	}
	var openBrowserCommands = []string{
		"wsl-open",           // for Windows Subsystem for Linux, see https://github.com/git-town/git-town/issues/1344
		"garcon-url-handler", // opens links in native browser on ChromeOS
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
		res, err := command.Run("which", browserCommand)
		if err == nil && res.OutputSanitized() != "" {
			return browserCommand
		}
	}
	return ""
}

// Open opens the default browser with the given URL.
// If no browser is found, prints the URL.
func Open(url string, shell command.Shell) {
	command := OpenBrowserCommand()
	if command == "" {
		fmt.Println("Please open in a browser: " + url)
		return
	}
	_, err := shell.Run(command, url)
	if err != nil {
		fmt.Println("Please open in a browser: " + url)
	}
}
