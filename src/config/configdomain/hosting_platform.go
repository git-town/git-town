package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/messages"
)

// HostingPlatform defines legal values for the "git-town.code-hosting-platform" config setting.
type HostingPlatform string

func (self HostingPlatform) String() string { return string(self) }

const (
	HostingPlatformBitbucket = HostingPlatform("bitbucket")
	HostingPlatformGitHub    = HostingPlatform("github")
	HostingPlatformGitLab    = HostingPlatform("gitlab")
	HostingPlatformGitea     = HostingPlatform("gitea")
	HostingPlatformNone      = HostingPlatform("")
)

// NewHostingPlatform provides the HostingService enum matching the given text.
func NewHostingPlatform(platformName string) (HostingPlatform, error) {
	text := strings.ToLower(platformName)
	for _, hostingService := range hostingPlatforms() {
		if strings.ToLower(text) == hostingService.String() {
			return hostingService, nil
		}
	}
	return HostingPlatformNone, fmt.Errorf(messages.HostingServiceUnknown, text)
}

// NewHostingPlatformRef provides the HostingService enum matching the given text.
func NewHostingPlatformRef(platformName string) (*HostingPlatform, error) {
	result, err := NewHostingPlatform(platformName)
	return &result, err
}

// hostingPlatforms provides all legal values for HostingService.
func hostingPlatforms() []HostingPlatform {
	return []HostingPlatform{
		HostingPlatformNone,
		HostingPlatformBitbucket,
		HostingPlatformGitHub,
		HostingPlatformGitLab,
		HostingPlatformGitea,
	}
}
