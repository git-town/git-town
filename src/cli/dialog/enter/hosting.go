package enter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

const enterHostingPlatformHelp = `
Git Town needs to know on which code hosting platform you host your code.
This allows it to open the right browser URLs and talk to the right API endpoints.
Only change this setting if the auto-detection does not work for you.

`

// EnterMainBranch lets the user select a new main branch for this repo.
func HostingPlatform(existingValue configdomain.HostingPlatform, inputs dialogcomponents.TestInput) (configdomain.HostingPlatform, bool, error) {
	entries := []hostingPlatformEntry{
		hostingPlatformAutoDetect,
		hostingPlatformBitBucket,
		hostingPlatformGitea,
		hostingPlatformGitHub,
		hostingPlatformGitLab,
	}
	cursor := indexOfHostingPlatform(existingValue, entries)
	newValue, aborted, err := dialogcomponents.RadioList(entries, cursor, enterHostingPlatformHelp, inputs)
	fmt.Printf("Code hosting: %s\n", dialogcomponents.FormattedSelection(newValue.String(), aborted))
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
