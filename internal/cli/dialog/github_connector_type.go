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
	githubConnectorTypeTitle = `GitHub connector`
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

func GithubConnectorType(args Args[forgedomain.GithubConnectorType]) (Option[forgedomain.GithubConnectorType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.GithubConnectorType]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.GithubConnectorType]]{
			Data: None[forgedomain.GithubConnectorType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.GithubConnectorType]]{
		{
			Data: Some(forgedomain.GithubConnectorTypeAPI),
			Text: "API token",
		},
		{
			Data: Some(forgedomain.GithubConnectorTypeGh),
			Text: "gh tool",
		},
	}...)
	cursor := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, githubConnectorTypeTitle, githubConnectorTypeHelp, args.Inputs, "github-connector")
	fmt.Printf(messages.GithubConnectorTypeResult, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
