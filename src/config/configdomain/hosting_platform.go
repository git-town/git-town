package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v14/src/gohacks"
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
)

// NewHostingPlatform provides the HostingPlatform enum matching the given text.
func NewHostingPlatform(platformName string) (HostingPlatform, error) {
	text := strings.ToLower(platformName)
	for _, hostingPlatform := range hostingPlatforms() {
		if strings.ToLower(text) == hostingPlatform.String() {
			return hostingPlatform, nil
		}
	}
	return HostingPlatformGitHub, fmt.Errorf(messages.HostingPlatformUnknown, text)
}

// NewHostingPlatformOption provides the HostingPlatform enum matching the given text.
func NewHostingPlatformOption(platformName string) (gohacks.Option[HostingPlatform], error) {
	platform, err := NewHostingPlatform(platformName)
	if err != nil {
		return gohacks.NewOptionNone[HostingPlatform](), err
	}
	return gohacks.NewOption(platform), nil
}

// hostingPlatforms provides all legal values for HostingPlatform.
func hostingPlatforms() []HostingPlatform {
	return []HostingPlatform{
		HostingPlatformBitbucket,
		HostingPlatformGitHub,
		HostingPlatformGitLab,
		HostingPlatformGitea,
	}
}
