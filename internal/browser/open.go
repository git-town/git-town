// Package browser allows interacting with the default browser on the user's machine.
package browser

import (
	"fmt"
	"runtime"

	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	"github.com/git-town/git-town/v22/internal/filesystem"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Open opens a new window/tab in the default browser with the given URL.
// If no browser is found, it prints the URL.
func Open(url string, frontend subshelldomain.Runner, config Option[browserdomain.Browser]) {
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
func OpenBrowserCommand(config Option[browserdomain.Browser]) Option[string] {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return Some("start")
	}
	browserCommands, useBrowser := browserCommandsToUse(config).Get()
	if !useBrowser {
		return None[string]()
	}
	return filesystem.FirstExistingExecutable(browserCommands)
}

// browserCommandsToUse provides the browser commands to use based on the config.
// A None result means that the user wants to use no browser.
func browserCommandsToUse(browserConfig Option[browserdomain.Browser]) Option[[]string] {
	userBrowser, hasUserBrowser := browserConfig.Get()
	if !hasUserBrowser {
		return Some(defaultBrowserCommands())
	}
	if userBrowser.NoBrowser() {
		return None[[]string]()
	}
	return Some(append([]string{userBrowser.String()}, defaultBrowserCommands()...))
}

// defaultBrowserCommands provides the default browser commands Git Town knows about.
func defaultBrowserCommands() []string {
	return []string{
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
}
