package config

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/messages"
)

// HostingService defines legal values for the "git-town.code-hosting-driver" config setting.
type HostingService string

const (
	HostingServiceBitbucket HostingService = "bitbucket"
	HostingServiceGitHub    HostingService = "github"
	HostingServiceGitLab    HostingService = "gitlab"
	HostingServiceGitea     HostingService = "gitea"
	HostingServiceNone      HostingService = ""
)

// NewHostingService provides the HostingService enum matching the given text.
func NewHostingService(text string) (HostingService, error) {
	text = strings.ToLower(text)
	for _, hostingService := range hostingServices() {
		if string(hostingService) == text {
			return hostingService, nil
		}
	}
	return HostingServiceNone, fmt.Errorf(messages.HostingServiceUnknown, text)
}

// hostingServices provides all legal values for HostingService.
func hostingServices() []HostingService {
	return []HostingService{
		HostingServiceNone,
		HostingServiceBitbucket,
		HostingServiceGitHub,
		HostingServiceGitLab,
		HostingServiceGitea,
	}
}
