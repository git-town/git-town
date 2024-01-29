package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

const enterOriginHostnameHelp = `
When using SSH identities, define the hostname of your source code repository.

Only change this setting if the auto-detection does not work for you.

`

func OriginHostname(existingValue configdomain.HostingOriginHostname, inputs components.TestInput) (configdomain.HostingPlatform, bool, error) {
	entries := []hostingPlatformEntry{
		hostingPlatformAutoDetect,
		hostingPlatformBitBucket,
		hostingPlatformGitea,
		hostingPlatformGitHub,
		hostingPlatformGitLab,
	}
	cursor := indexOfHostingPlatform(existingValue, entries)
	newValue, aborted, err := components.RadioList(entries, cursor, enterHostingPlatformHelp, inputs)
	fmt.Printf("Code hosting: %s\n", components.FormattedSelection(newValue.String(), aborted))
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
