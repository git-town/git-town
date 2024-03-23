package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/gohacks/stringers"
	"github.com/git-town/git-town/v13/src/messages"
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
	entries := []hostingPlatformEntry{
		hostingPlatformAutoDetect,
		hostingPlatformBitBucket,
		hostingPlatformGitea,
		hostingPlatformGitHub,
		hostingPlatformGitLab,
	}
	cursor := indexOfHostingPlatform(existingValue, entries)
	newValue, aborted, err := components.RadioList(entries, cursor, hostingPlatformTitle, HostingPlatformHelp, inputs)
	fmt.Printf(messages.CodeHosting, components.FormattedSelection(newValue.String(), aborted))
	return newValue.HostingPlatform(), aborted, err
}

type hostingPlatformEntry string

const (
	hostingPlatformAutoDetect hostingPlatformEntry = "auto-detect"
	hostingPlatformBitBucket  hostingPlatformEntry = "BitBucket"
	hostingPlatformGitea      hostingPlatformEntry = "Gitea"
	hostingPlatformGitHub     hostingPlatformEntry = "Github"
	hostingPlatformGitLab     hostingPlatformEntry = "GitLab"
)

func (self hostingPlatformEntry) HostingPlatform() configdomain.HostingPlatform {
	switch self {
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
	panic("unknown hosting platform: " + self)
}

func (self hostingPlatformEntry) String() string {
	return string(self)
}

func indexOfHostingPlatform(hostingPlatform configdomain.HostingPlatform, entries []hostingPlatformEntry) int {
	entry := newHostingPlatformEntry(hostingPlatform)
	return stringers.IndexOrStart(entries, entry)
}

func newHostingPlatformEntry(hosting configdomain.HostingPlatform) hostingPlatformEntry {
	switch hosting {
	case configdomain.HostingPlatformNone:
		return hostingPlatformAutoDetect
	case configdomain.HostingPlatformBitbucket:
		return hostingPlatformBitBucket
	case configdomain.HostingPlatformGitea:
		return hostingPlatformGitea
	case configdomain.HostingPlatformGitHub:
		return hostingPlatformGitHub
	case configdomain.HostingPlatformGitLab:
		return hostingPlatformGitLab
	}
	panic("unknown hosting: " + hosting)
}
