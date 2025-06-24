package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	tokenScopeTitle = `API token scope`
	tokenScopeHelp  = `
Do you want to use the API token you just entered
just for this Git repo or for all Git repos on this machine?

`
)

// GitHubToken lets the user enter the GitHub API token.
func TokenScope(oldValue configdomain.ConfigScope, inputs components.TestInput) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	entries := list.Entries[configdomain.ConfigScope]{
		{
			Data: configdomain.ConfigScopeGlobal,
			Text: "globally for all Git repositories on my machine",
		},
		{
			Data: configdomain.ConfigScopeLocal,
			Text: "locally only for the this Git repository",
		},
	}
	defaultPos := entries.IndexOf(oldValue)
	selection, exit, err := components.RadioList(entries, defaultPos, tokenScopeTitle, tokenScopeHelp, inputs)
	if err != nil || exit {
		return configdomain.ConfigScopeLocal, exit, err
	}
	fmt.Printf(messages.ForgeAPITokenLocation, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
