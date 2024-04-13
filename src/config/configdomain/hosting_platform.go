package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v14/src/messages"
)

// HostingPlatform defines legal values for the "git-town.hosting-platform" config setting.
type HostingPlatform string

func (self HostingPlatform) String() string { return string(self) }

const (
	HostingPlatformBitbucket = HostingPlatform("bitbucket")
	HostingPlatformGitHub    = HostingPlatform("github")
	HostingPlatformGitLab    = HostingPlatform("gitlab")
	HostingPlatformGitea     = HostingPlatform("gitea")
	HostingPlatformNone      = HostingPlatform("") // no hosting or auto-detect
)

// NewHostingPlatform provides the HostingPlatform enum matching the given text.
func NewHostingPlatform(platformName string) (HostingPlatform, error) {
	text := strings.ToLower(platformName)
	for _, hostingPlatform := range hostingPlatforms() {
		if strings.ToLower(text) == hostingPlatform.String() {
			return hostingPlatform, nil
		}
	}
	return HostingPlatformNone, fmt.Errorf(messages.HostingPlatformUnknown, text)
}

// NewHostingPlatformRef provides the HostingPlatform enum matching the given text.
func NewHostingPlatformRef(platformName string) (*HostingPlatform, error) {
	result, err := NewHostingPlatform(platformName)
	return &result, err
}

// hostingPlatforms provides all legal values for HostingPlatform.
func hostingPlatforms() []HostingPlatform {
	return []HostingPlatform{
		HostingPlatformNone,
		HostingPlatformBitbucket,
		HostingPlatformGitHub,
		HostingPlatformGitLab,
		HostingPlatformGitea,
	}
}
