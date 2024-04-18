package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
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

func HostingPlatform(existingValue configdomain.HostingPlatform, inputs components.TestInput) (configdomain.HostingPlatform, bool, error) {
	entries := hostingPlatformEntries{
		hostingPlatformAutoDetect,
		hostingPlatformBitBucket,
		hostingPlatformGitea,
		hostingPlatformGitHub,
		hostingPlatformGitLab,
	}
	cursor := entries.IndexOfHostingPlatformOrStart(existingValue)
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

func (hpe hostingPlatformEntry) HostingPlatform() configdomain.HostingPlatform {
	switch hpe {
	case hostingPlatformAutoDetect:
		return configdomain.HostingPlatformNone
	case hostingPlatformBitBucket:
		return configdomain.HostingPlatformBitbucket
	case hostingPlatformGitea:
		return configdomain.HostingPlatformGitea
	case hostingPlatformGitHub:
		return configdomain.HostingPlatformGitHub
	case hostingPlatformGitLab:
		return configdomain.HostingPlatformGitLab
	}
	panic("unknown hosting platform: " + hpe)
}

func (hpe hostingPlatformEntry) String() string {
	return string(hpe)
}

type hostingPlatformEntries []hostingPlatformEntry

func (hpes hostingPlatformEntries) IndexOfHostingPlatformOrStart(needle configdomain.HostingPlatform) int {
	for h, hostingPlatformEntry := range hpes {
		if hostingPlatformEntry.HostingPlatform() == needle {
			return h
		}
	}
	return 0
}
