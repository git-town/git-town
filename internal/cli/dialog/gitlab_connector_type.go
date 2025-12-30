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
	gitLabConnectorTypeTitle = `GitLab connector`
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

func GitLabConnectorType(args Args[forgedomain.GitLabConnectorType]) (Option[forgedomain.GitLabConnectorType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.GitLabConnectorType]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.GitLabConnectorType]]{
			Data: None[forgedomain.GitLabConnectorType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.GitLabConnectorType]]{
		{
			Data: Some(forgedomain.GitLabConnectorTypeAPI),
			Text: "API token",
		},
		{
			Data: Some(forgedomain.GitLabConnectorTypeGlab),
			Text: "glab tool",
		},
	}...)
	cursor := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, gitLabConnectorTypeTitle, gitLabConnectorTypeHelp, args.Inputs, "gitlab-connector")
	fmt.Printf(messages.GitLabConnectorTypeResult, dialogcomponents.FormattedSelection(selection.GetOrZero().String(), exit))
	return selection, exit, err
}
