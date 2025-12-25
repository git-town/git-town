package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

const (
	tokenScopeTitle = `API token scope`
	tokenScopeHelp  = `
Do you want to use the API token
you just entered just for this Git repo
or for all Git repos on this machine?

`
)

// TokenScope lets the user enter the GitHub API token.
func TokenScope(oldValue configdomain.ConfigScope, inputs dialogcomponents.Inputs) (configdomain.ConfigScope, dialogdomain.Exit, error) {
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
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, tokenScopeTitle, tokenScopeHelp, inputs, "token-scope")
	fmt.Printf(messages.ForgeAPITokenLocation, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
