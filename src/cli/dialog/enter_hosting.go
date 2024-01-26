package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

const enterHostingPlatformHelp = `
Git Town needs to know on which code hosting platform you host your code.
This allows it to open the right browser URLs and talk to the right API endpoints.
Only change this setting if the auto-detection does not work for you.

`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterHostingPlatform(existingValue configdomain.HostingPlatform, inputs TestInput) (configdomain.HostingPlatform, bool, error) {
	entries := []hostingPlatformEntry{
		hostingPlatformAutoDetect,
		hostingPlatformBitBucket,
		hostingPlatformGitea,
		hostingPlatformGitHub,
		hostingPlatformGitLab,
	}
	cursor := indexOfHostingPlatform(existingValue, entries)
	newValue, aborted, err := radioList(entries, cursor, enterHostingPlatformHelp, inputs)
	fmt.Printf("Code hosting: %s\n", formattedSelection(newValue.String(), aborted))
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
		return configdomain.HostingNone
	case hostingPlatformBitBucket:
		return configdomain.HostingBitbucket
	case hostingPlatformGitea:
		return configdomain.HostingGitea
	case hostingPlatformGitHub:
		return configdomain.HostingGitHub
	case hostingPlatformGitLab:
		return configdomain.HostingGitLab
	}
	panic("unknown hosting platform: " + self)
}

func (self hostingPlatformEntry) String() string {
	return string(self)
}

func newHostingPlatformEntry(hosting configdomain.HostingPlatform) hostingPlatformEntry {
	switch hosting {
	case configdomain.HostingNone:
		return hostingPlatformAutoDetect
	case configdomain.HostingBitbucket:
		return hostingPlatformBitBucket
	case configdomain.HostingGitea:
		return hostingPlatformGitea
	case configdomain.HostingGitHub:
		return hostingPlatformGitHub
	case configdomain.HostingGitLab:
		return hostingPlatformGitLab
	}
	panic("unknown hosting: " + hosting)
}

func indexOfHostingPlatform(hostingPlatform configdomain.HostingPlatform, entries []hostingPlatformEntry) int {
	entry := newHostingPlatformEntry(hostingPlatform)
	return stringers.IndexOrStart(entries, entry)
}
