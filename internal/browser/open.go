// Package browser allows interacting with the default browser on the user's machine.
package browser

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
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
func OpenBrowserCommand(customBrowserSetting Option[browserdomain.Browser]) Option[string] {
	customBrowser, hasCustomBrowser := customBrowserSetting.Get()
	if !hasCustomBrowser {
		return defaultBrowserCommand()
	}
	customBrowserCmd, useBrowser := customBrowser.Get()
	if !useBrowser {
		return None[string]()
	}
	return Some(customBrowserCmd)
}
