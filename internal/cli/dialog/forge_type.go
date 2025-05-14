package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

const (
	hostingPlatformTitle = `Forge type`
	HostingPlatformHelp  = `
Git Town uses this setting to open browser URLs
and interact with your code hosting provider's API.

Only change this if your forge is hosted at a custom URL.

`
)

func ForgeType(existingValue Option[forgedomain.ForgeType], inputs components.TestInput) (Option[forgedomain.ForgeType], bool, error) {
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
	newValue, aborted, err := components.RadioList(entries, cursor, hostingPlatformTitle, HostingPlatformHelp, inputs)
	fmt.Printf(messages.Forge, components.FormattedSelection(newValue.GetOrElse("auto-detect").String(), aborted))
	return newValue, aborted, err
}
