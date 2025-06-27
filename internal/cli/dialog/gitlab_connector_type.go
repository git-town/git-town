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
	gitLabConnectorTypeTitle = `GitLab connector type`
	gitLabConnectorTypeHelp  = `
Git Town supports two ways to connect to GitLab:

1. GitLab API:
   Git Town talks directly with the GitLab API.
	 This uses an access token for your account
	 that you provide on the next screen.
2. GitLab's "glab" application:
   This doesn't require an access token
	 but you need to install and configure the glab tool.

`
)

func GitLabConnectorType(existing Option[forgedomain.GitLabConnectorType], inputs components.TestInput) (forgedomain.GitLabConnectorType, dialogdomain.Exit, error) {
	entries := list.Entries[forgedomain.GitLabConnectorType]{
		{
			Data: forgedomain.GitLabConnectorTypeAPI,
			Text: "API token",
		},
		{
			Data: forgedomain.GitLabConnectorTypeGlab,
			Text: "glab tool",
		},
	}
	defaultPos := 0
	if existingValue, hasExisting := existing.Get(); hasExisting {
		defaultPos = entries.IndexOf(existingValue)
	}
	selection, exit, err := components.RadioList(entries, defaultPos, gitLabConnectorTypeTitle, gitLabConnectorTypeHelp, inputs)
	if err != nil || exit {
		return forgedomain.GitLabConnectorTypeAPI, exit, err
	}
	fmt.Printf(messages.GitLabConnectorType, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
