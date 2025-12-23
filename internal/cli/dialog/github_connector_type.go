package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	gitHubConnectorTypeTitle = `GitHub connector`
	gitHubConnectorTypeHelp  = `
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

func GitHubConnectorType(args Args[forgedomain.GitHubConnectorType]) (Option[forgedomain.GitHubConnectorType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.GitHubConnectorType]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.GitHubConnectorType]]{
			Data: None[forgedomain.GitHubConnectorType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.GitHubConnectorType]]{
		{
			Data: Some(forgedomain.GitHubConnectorTypeAPI),
			Text: "API token",
		},
		{
			Data: Some(forgedomain.GitHubConnectorTypeGh),
			Text: "gh tool",
		},
	}...)
	cursor := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, gitHubConnectorTypeTitle, gitHubConnectorTypeHelp, args.Inputs, "github-connector")
	fmt.Printf(messages.GitHubConnectorTypeResult, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
