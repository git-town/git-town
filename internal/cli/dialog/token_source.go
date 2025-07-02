package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	tokenSourceTitle = `%s token source`
)

// GitHubToken lets the user enter the GitHub API token.
func TokenSource(forgeType forgedomain.ForgeType, oldValue forgedomain.TokenSource, inputs components.TestInput) (forgedomain.TokenSource, dialogdomain.Exit, error) {
	entries := list.Entries[forgedomain.TokenSource]{
		{
			Data: forgedomain.TokenSourceManual,
			Text: "enter the token directly",
		},
		{
			Data: forgedomain.TokenSourceScript,
			Text: "provide a script that outputs the token to Git Town if needed",
		},
	}
	defaultPos := entries.IndexOf(oldValue)
	selection, exit, err := components.RadioList(entries, defaultPos, fmt.Sprintf(tokenSourceTitle, forgeType), "", inputs)
	fmt.Printf(messages.ForgeAPITokenSource, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
