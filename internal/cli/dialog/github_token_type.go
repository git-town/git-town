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

func GitHubTokenType(existing Option[forgedomain.GitHubTokenType], inputs dialogcomponents.Inputs) (forgedomain.GitHubTokenType, dialogdomain.Exit, error) {
	entries := list.Entries[forgedomain.GitHubTokenType]{
		{
			Data: forgedomain.GitHubTokenTypeEnter,
			Text: "enter the token directly",
		},
		{
			Data: forgedomain.GitHubTokenTypeCLI,
			Text: "enter a shell call that provides the token",
		},
	}
	defaultPos := 0
	if existingValue, hasExisting := existing.Get(); hasExisting {
		defaultPos = entries.IndexOf(existingValue)
	}
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, gitHubTokenTypeTitle, gitHubTokenTypeHelp, inputs, "github-token-type")
	if err != nil || exit {
		return forgedomain.GitHubTokenTypeEnter, exit, err
	}
	fmt.Printf(messages.GitHubTokenType, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
