package config

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/messages"
)

// Hosting defines legal values for the "git-town.code-hosting-driver" config setting.
type Hosting struct {
	name string
}

func (h Hosting) String() string { return h.name }

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
	for _, hostingService := range hostingServices() {
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
