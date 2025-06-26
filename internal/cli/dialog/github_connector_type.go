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

1. GitHub API:
   Git Town talks directly with the GitHub API.
	 This uses an access token for your account
	 that you provide on the next screen.
2. GitHub's "gh" application:
   This doesn't require an access token
	 but you need to install and configure the gh tool.

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
