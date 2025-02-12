package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v18/internal/messages"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// HostingPlatform defines legal values for the "git-town.hosting-platform" config setting.
type HostingPlatform string

func (self HostingPlatform) String() string { return string(self) }

const (
	HostingPlatformBitbucket           = HostingPlatform("bitbucket")
	HostingPlatformBitbucketDatacenter = HostingPlatform("bitbucket-datacenter")
	HostingPlatformGitHub              = HostingPlatform("github")
	HostingPlatformGitLab              = HostingPlatform("gitlab")
	HostingPlatformGitea               = HostingPlatform("gitea")
)

// ParseHostingPlatform provides the HostingPlatform enum matching the given text.
func ParseHostingPlatform(platformName string) (Option[HostingPlatform], error) {
	if platformName == "" {
		return None[HostingPlatform](), nil
	}
	platformNameLower := strings.ToLower(platformName)
	for _, hostingPlatform := range hostingPlatforms() {
		if platformNameLower == hostingPlatform.String() {
			return Some(hostingPlatform), nil
		}
	}
	return None[HostingPlatform](), fmt.Errorf(messages.HostingPlatformUnknown, platformName)
}

// hostingPlatforms provides all legal values for HostingPlatform.
func hostingPlatforms() []HostingPlatform {
	return []HostingPlatform{
		HostingPlatformBitbucket,
		HostingPlatformBitbucketDatacenter,
		HostingPlatformGitHub,
		HostingPlatformGitLab,
		HostingPlatformGitea,
	}
}
