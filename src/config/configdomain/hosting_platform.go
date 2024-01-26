package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/messages"
)

// HostingPlatform defines legal values for the "git-town.code-hosting-platform" config setting.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type HostingPlatform string

func (self HostingPlatform) String() string { return string(self) }

const (
	HostingBitbucket = HostingPlatform("bitbucket")
	HostingGitHub    = HostingPlatform("github")
	HostingGitLab    = HostingPlatform("gitlab")
	HostingGitea     = HostingPlatform("gitea")
	HostingNone      = HostingPlatform("")
)

// NewHosting provides the HostingService enum matching the given text.
func NewHosting(platformName string) (HostingPlatform, error) {
	text := strings.ToLower(platformName)
	for _, hostingService := range hostings() {
		if strings.ToLower(text) == hostingService.String() {
			return hostingService, nil
		}
	}
	return HostingNone, fmt.Errorf(messages.HostingServiceUnknown, text)
}

// NewHostingRef provides the HostingService enum matching the given text.
func NewHostingRef(platformName string) (*HostingPlatform, error) {
	result, err := NewHosting(platformName)
	return &result, err
}

// hostings provides all legal values for HostingService.
func hostings() []HostingPlatform {
	return []HostingPlatform{
		HostingNone,
		HostingBitbucket,
		HostingGitHub,
		HostingGitLab,
		HostingGitea,
	}
}
