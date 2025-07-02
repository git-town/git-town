package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	tokenSourceTitle = `%s token source`
)

// GitHubToken lets the user enter the GitHub API token.
func TokenSource(forgeType forgedomain.ForgeType, oldValue Option[forgedomain.TokenSource], inputs components.TestInput) (forgedomain.TokenSource, dialogdomain.Exit, error) {
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
	defaultPos := 0
	if existing, has := oldValue.Get(); has {
		entries.IndexOf(existing)
	}
	selection, exit, err := components.RadioList(entries, defaultPos, fmt.Sprintf(tokenSourceTitle, forgeType), "", inputs)
	fmt.Printf(messages.ForgeAPITokenSource, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
