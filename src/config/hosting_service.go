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

const (
	HostingBitbucket Hosting = "bitbucket"
	HostingGitHub    Hosting = "github"
	HostingGitLab    Hosting = "gitlab"
	HostingGitea     Hosting = "gitea"
	HostingNone      Hosting = ""
)

// NewHosting provides the HostingService enum matching the given text.
func NewHosting(text string) (Hosting, error) {
	text = strings.ToLower(text)
	for _, hostingService := range hostings() {
		if string(hostingService) == text {
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
