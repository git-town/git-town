package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

const (
	hostingPlatformTitle = `Hosting platform`
	HostingPlatformHelp  = `
Knowing the type of code hosting platform allows Git Town
to open browser URLs and talk to the code hosting API.
Most people can leave this on "auto-detect".
Only change this if your code hosting server uses as custom URL.

`
)

func HostingPlatform(existingValue Option[configdomain.HostingPlatform], inputs components.TestInput) (Option[configdomain.HostingPlatform], bool, error) {
	entries := list.Entries[Option[configdomain.HostingPlatform]]{
		{
			Data:    None[configdomain.HostingPlatform](),
			Enabled: true,
			Text:    "auto-detect",
		},
		{
			Data:    Some(configdomain.HostingPlatformBitbucket),
			Enabled: true,
			Text:    "BitBucket",
		},
		{
			Data:    Some(configdomain.HostingPlatformBitbucketDatacenter),
			Enabled: true,
			Text:    "BitBucket-Datacenter",
		},
		{
			Data:    Some(configdomain.HostingPlatformGitea),
			Enabled: true,
			Text:    "Gitea",
		},
		{
			Data:    Some(configdomain.HostingPlatformGitHub),
			Enabled: true,
			Text:    "GitHub",
		},
		{
			Data:    Some(configdomain.HostingPlatformGitLab),
			Enabled: true,
			Text:    "GitLab",
		},
	}
	cursor := entries.IndexOf(existingValue)
	newValue, aborted, err := components.RadioList(entries, cursor, hostingPlatformTitle, HostingPlatformHelp, inputs)
	fmt.Printf(messages.CodeHosting, components.FormattedSelection(newValue.GetOrElse("auto-detect").String(), aborted))
	return newValue, aborted, err
}
