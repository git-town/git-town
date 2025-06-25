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
	githubConnectorTypeTitle = `GitHub connector type`
	githubConnectorTypeHelp  = `
Git Town supports two ways to connect to GitHub:
1. via the GitHub API
2. via GitHub's "gh" tool

Option 1 requires you to create an API token and enter it here.
Option 2 requires you to install the "gh" tool.

`
)

func GitHubConnectorType(existing Option[forgedomain.GitHubConnectorType], inputs components.TestInput) (forgedomain.GitHubConnectorType, dialogdomain.Exit, error) {
	entries := list.Entries[forgedomain.GitHubConnectorType]{
		{
			Data: forgedomain.GitHubConnectorTypeAPI,
			Text: "GitHub API token",
		},
		{
			Data: forgedomain.GitHubConnectorTypeGh,
			Text: "gh tool",
		},
	}
	defaultPos := 0
	if existingValue, hasExisting := existing.Get(); hasExisting {
		defaultPos = entries.IndexOf(existingValue)
	}
	selection, exit, err := components.RadioList(entries, defaultPos, githubConnectorTypeTitle, githubConnectorTypeHelp, inputs)
	if err != nil || exit {
		return forgedomain.GitHubConnectorTypeAPI, exit, err
	}
	fmt.Printf(messages.GitHubConnectorType, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
