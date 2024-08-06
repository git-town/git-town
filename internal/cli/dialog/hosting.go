package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/messages"
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
	entries := hostingPlatformEntries{
		hostingPlatformAutoDetect,
		hostingPlatformBitBucket,
		hostingPlatformGitea,
		hostingPlatformGitHub,
		hostingPlatformGitLab,
	}
	cursor := entries.IndexOfHostingPlatformOrStart(existingValue.GetOrDefault())
	newValue, aborted, err := components.RadioList(list.NewEntries(entries...), cursor, hostingPlatformTitle, HostingPlatformHelp, inputs)
	fmt.Printf(messages.CodeHosting, components.FormattedSelection(newValue.String(), aborted))
	return newValue.HostingPlatform(), aborted, err
}

type hostingPlatformEntry string

const (
	hostingPlatformAutoDetect hostingPlatformEntry = "auto-detect"
	hostingPlatformBitBucket  hostingPlatformEntry = "BitBucket"
	hostingPlatformGitea      hostingPlatformEntry = "Gitea"
	hostingPlatformGitHub     hostingPlatformEntry = "GitHub"
	hostingPlatformGitLab     hostingPlatformEntry = "GitLab"
)

func (entry hostingPlatformEntry) HostingPlatform() Option[configdomain.HostingPlatform] {
	switch entry {
	case hostingPlatformAutoDetect, "":
		return None[configdomain.HostingPlatform]()
	case hostingPlatformBitBucket:
		return Some(configdomain.HostingPlatformBitbucket)
	case hostingPlatformGitea:
		return Some(configdomain.HostingPlatformGitea)
	case hostingPlatformGitHub:
		return Some(configdomain.HostingPlatformGitHub)
	case hostingPlatformGitLab:
		return Some(configdomain.HostingPlatformGitLab)
	}
	panic(fmt.Sprintf(messages.HostingPlatformUnknown, entry))
}

func (entry hostingPlatformEntry) String() string {
	return string(entry)
}

type hostingPlatformEntries []hostingPlatformEntry

func (entries hostingPlatformEntries) IndexOfHostingPlatformOrStart(needle configdomain.HostingPlatform) int {
	for h, entry := range entries {
		if value, has := entry.HostingPlatform().Get(); has && value == needle {
			return h
		}
	}
	return 0
}
