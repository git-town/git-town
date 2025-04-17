package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

const (
	hostingPlatformTitle = `Forge type`
	HostingPlatformHelp  = `
Knowing the type of forge allows Git Town
to open browser URLs and talk to the forge API.
Most people can leave this on "auto-detect".
Only change this if your forge uses as custom URL.

`
)

func ForgeType(existingValue Option[configdomain.ForgeType], inputs components.TestInput) (Option[configdomain.ForgeType], bool, error) {
	entries := list.Entries[Option[configdomain.ForgeType]]{
		{
			Data: None[configdomain.ForgeType](),
			Text: "auto-detect",
		},
		{
			Data: Some(configdomain.ForgeTypeBitbucket),
			Text: "Bitbucket",
		},
		{
			Data: Some(configdomain.ForgeTypeBitbucketDatacenter),
			Text: "Bitbucket Data Center",
		},
		{
			Data: Some(configdomain.ForgeTypeCodeberg),
			Text: "Codeberg",
		},
		{
			Data: Some(configdomain.ForgeTypeGitea),
			Text: "Gitea",
		},
		{
			Data: Some(configdomain.ForgeTypeGitHub),
			Text: "GitHub",
		},
		{
			Data: Some(configdomain.ForgeTypeGitLab),
			Text: "GitLab",
		},
	}
	cursor := entries.IndexOfFunc(existingValue, func(optA, optB Option[configdomain.ForgeType]) bool {
		return optA.Equal(optB)
	})
	newValue, aborted, err := components.RadioList(entries, cursor, hostingPlatformTitle, HostingPlatformHelp, inputs)
	fmt.Printf(messages.Forge, components.FormattedSelection(newValue.GetOrElse("auto-detect").String(), aborted))
	return newValue, aborted, err
}
