package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/messages"
)

// Hosting defines legal values for the "git-town.code-hosting-platform" config setting.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type Hosting string

func (self Hosting) String() string { return string(self) }

const (
	HostingBitbucket = Hosting("bitbucket")
	HostingGitHub    = Hosting("github")
	HostingGitLab    = Hosting("gitlab")
	HostingGitea     = Hosting("gitea")
	HostingNone      = Hosting("")
)

// NewHosting provides the HostingService enum matching the given text.
func NewHosting(platformName string) (Hosting, error) {
	text := strings.ToLower(platformName)
	for _, hostingService := range hostings() {
		if strings.ToLower(text) == hostingService.String() {
			return hostingService, nil
		}
	}
	return HostingNone, fmt.Errorf(messages.HostingServiceUnknown, text)
}

// NewHostingRef provides the HostingService enum matching the given text.
func NewHostingRef(platformName string) (*Hosting, error) {
	result, err := NewHosting(platformName)
	return &result, err
}

// hostings provides all legal values for HostingService.
func hostings() []Hosting {
	return []Hosting{
		HostingNone,
		HostingBitbucket,
		HostingGitHub,
		HostingGitLab,
		HostingGitea,
	}
}
