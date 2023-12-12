package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/messages"
)

// Hosting defines legal values for the "git-town.code-hosting-platform" config setting.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type Hosting struct {
	name string
}

func (self Hosting) String() string { return self.name }

var (
	HostingBitbucket = Hosting{"bitbucket"} //nolint:gochecknoglobals
	HostingGitHub    = Hosting{"github"}    //nolint:gochecknoglobals
	HostingGitLab    = Hosting{"gitlab"}    //nolint:gochecknoglobals
	HostingGitea     = Hosting{"gitea"}     //nolint:gochecknoglobals
	HostingNone      = Hosting{""}          //nolint:gochecknoglobals
)

// NewHosting provides the HostingService enum matching the given text.
func NewHosting(text string) (Hosting, error) {
	text = strings.ToLower(text)
	for _, hostingService := range hostings() {
		if hostingService.name == text {
			return hostingService, nil
		}
	}
	return HostingNone, fmt.Errorf(messages.HostingServiceUnknown, text)
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
