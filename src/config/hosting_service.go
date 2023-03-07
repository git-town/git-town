package config

import "fmt"

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
	for _, hostingService := range hostingServices() {
		if string(hostingService) == text {
			return hostingService, nil
		}
	}
	return HostingServiceNone, fmt.Errorf("unknown alias type: %q", text)
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
