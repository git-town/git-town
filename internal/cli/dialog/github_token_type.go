package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	gitHubTokenTypeTitle = `GitHub token input`
	gitHubTokenTypeHelp  = `
Git Town supports multiple ways to enter the GitHub token:

`
)

func GitHubTokenType(args Args[forgedomain.GitHubTokenType]) (Option[forgedomain.GitHubTokenType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.GitHubTokenType]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.GitHubTokenType]]{
			Data: None[forgedomain.GitHubTokenType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.GitHubTokenType]]{
		{
			Data: Some(forgedomain.GitHubTokenTypeEnter),
			Text: "enter the token directly",
		},
		{
			Data: Some(forgedomain.GitHubTokenTypeCLI),
			Text: "enter a shell call that provides the token",
		},
	}...)
	cursor := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, gitHubTokenTypeTitle, gitHubTokenTypeHelp, args.Inputs, "github-token-type")
	if err != nil || exit {
		return None[forgedomain.GitHubTokenType](), exit, err
	}
	fmt.Printf(messages.GitHubTokenType, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, false, nil
}
