// Package browser allows interacting with the default browser on the user's machine.
package browser

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Open opens a new window/tab in the default browser with the given URL.
// If no browser is found, it prints the URL.
func Open(url string, frontend subshelldomain.Runner, config Option[configdomain.Browser]) {
	command, hasCommand := OpenBrowserCommand(config).Get()
	if !hasCommand {
		fmt.Printf(messages.BrowserOpen, url)
		return
	}
	if err := frontend.Run(command, url); err != nil {
		fmt.Printf(messages.BrowserOpen, url)
	}
}

// OpenBrowserCommand provides the console command to open the default browser.
func OpenBrowserCommand(config Option[configdomain.Browser]) Option[string] {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return Some("start")
	}
	openBrowserCommands := make([]string, 0, 11)
	if browser, hasBrowser := config.Get(); hasBrowser {
		if browser.NoBrowser() {
			return None[string]()
		}
		openBrowserCommands = append(openBrowserCommands, browser.String())
	}
	openBrowserCommands = append(openBrowserCommands,
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
	)
	for _, browserCommand := range openBrowserCommands {
		executable, err := exec.LookPath(browserCommand)
		if err == nil && len(executable) > 0 {
			return Some(browserCommand)
		}
	}
	return None[string]()
}
