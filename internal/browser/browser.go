// Package browser allows interacting with the default browser on the user's machine.
package browser

import (
	"fmt"
	"runtime"

	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/messages"
)

// Open opens a new window/tab in the default browser with the given URL.
// If no browser is found, it prints the URL.
func Open(url string, frontend frontendRunner, backend backendRunner) {
	command, hasCommand := OpenBrowserCommand(backend).Get()
	if !hasCommand {
		fmt.Printf(messages.BrowserOpen, url)
		return
	}
	err := frontend.Run(command, url)
	if err != nil {
		fmt.Printf(messages.BrowserOpen, url)
	}
}

// OpenBrowserCommand provides the console command to open the default browser.
func OpenBrowserCommand(runner backendRunner) Option[string] {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return Some("start")
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
		output, err := runner.Query("which", browserCommand)
		if err == nil && output != "" {
			return Some(browserCommand)
		}
	}
	return None[string]()
}

type frontendRunner interface {
	Run(executable string, args ...string) error
}

type backendRunner interface {
	Query(executable string, args ...string) (string, error)
}
