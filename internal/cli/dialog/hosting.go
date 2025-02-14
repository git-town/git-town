package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	. "github.com/git-town/git-town/v18/pkg/prelude"
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

func HostingPlatform(existingValue Option[configdomain.HostingPlatform], inputs components.TestInput) (Option[configdomain.HostingPlatform], bool, error) {
	entries := list.Entries[Option[configdomain.HostingPlatform]]{
		{
			Data: None[configdomain.HostingPlatform](),
			Text: "auto-detect",
		},
		{
			Data: Some(configdomain.HostingPlatformBitbucket),
			Text: "Bitbucket",
		},
		{
			Data: Some(configdomain.HostingPlatformBitbucketDatacenter),
			Text: "Bitbucket Data Center",
		},
		{
			Data: Some(configdomain.HostingPlatformGitea),
			Text: "Gitea",
		},
		{
			Data: Some(configdomain.HostingPlatformGitHub),
			Text: "GitHub",
		},
		{
			Data: Some(configdomain.HostingPlatformGitLab),
			Text: "GitLab",
		},
	}
	cursor := entries.IndexOfFunc(existingValue, func(optA, optB Option[configdomain.HostingPlatform]) bool {
		return optA.Equal(optB)
	})
	newValue, aborted, err := components.RadioList(entries, cursor, hostingPlatformTitle, HostingPlatformHelp, inputs)
	fmt.Printf(messages.Forge, components.FormattedSelection(newValue.GetOrElse("auto-detect").String(), aborted))
	return newValue, aborted, err
}
