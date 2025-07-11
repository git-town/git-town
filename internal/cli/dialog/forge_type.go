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
	forgeTypeTitle = `Forge type`
	forgeTypeHelp  = `
Git Town uses this setting
to open browser URLs
and interact with your forge's API.

Only change this if your forge
is hosted at a custom URL.

`
)

func ForgeType(existingValue Option[forgedomain.ForgeType], inputs dialogcomponents.TestInput) (Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.ForgeType]]{
		{
			Data: None[forgedomain.ForgeType](),
			Text: "auto-detect",
		},
		{
			Data: Some(forgedomain.ForgeTypeBitbucket),
			Text: "Bitbucket",
		},
		{
			Data: Some(forgedomain.ForgeTypeBitbucketDatacenter),
			Text: "Bitbucket Data Center",
		},
		{
			Data: Some(forgedomain.ForgeTypeCodeberg),
			Text: "Codeberg",
		},
		{
			Data: Some(forgedomain.ForgeTypeGitea),
			Text: "Gitea",
		},
		{
			Data: Some(forgedomain.ForgeTypeGitHub),
			Text: "GitHub",
		},
		{
			Data: Some(forgedomain.ForgeTypeGitLab),
			Text: "GitLab",
		},
	}
	cursor := entries.IndexOfFunc(existingValue, func(optA, optB Option[forgedomain.ForgeType]) bool {
		return optA.Equal(optB)
	})
	newValue, exit, err := dialogcomponents.RadioList(entries, cursor, forgeTypeTitle, forgeTypeHelp, inputs)
	fmt.Printf(messages.Forge, dialogcomponents.FormattedSelection(newValue.GetOrElse("auto-detect").String(), exit))
	return newValue, exit, err
}
