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
	forgeTypeTitle = `Forge type`
	forgeTypeHelp  = `
Git Town uses this setting
to open browser URLs
and interact with your forge's API.

Only change this if your forge
is hosted at a custom URL.

`
)

func ForgeType(args Args[forgedomain.ForgeType]) (Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.ForgeType]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.ForgeType]]{
			Data: None[forgedomain.ForgeType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.ForgeType]]{
		{
			Data: None[forgedomain.ForgeType](),
			Text: messages.AutoDetect,
		},
		{
			Data: Some(forgedomain.ForgeTypeAzureDevOps),
			Text: "Azure DevOps",
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
			Data: Some(forgedomain.ForgeTypeForgejo),
			Text: "Forgejo",
		},
		{
			Data: Some(forgedomain.ForgeTypeGitea),
			Text: "Gitea",
		},
		{
			Data: Some(forgedomain.ForgeTypeGithub),
			Text: "GitHub",
		},
		{
			Data: Some(forgedomain.ForgeTypeGitlab),
			Text: "GitLab",
		},
	}...)
	cursor := entries.IndexOfFunc(args.Local, func(optA, optB Option[forgedomain.ForgeType]) bool {
		return optA.Equal(optB)
	})
	newValue, exit, err := dialogcomponents.RadioList(entries, cursor, forgeTypeTitle, forgeTypeHelp, args.Inputs, "forge-type")
	if newValue.IsNone() {
		if args.Global.IsSome() {
			fmt.Printf(messages.Forge, dialogcomponents.FormattedOption(newValue, true, exit))
		} else {
			fmt.Printf(messages.Forge, dialogcomponents.FormattedSelection(messages.AutoDetect, exit))
		}
	}
	return newValue, exit, err
}
