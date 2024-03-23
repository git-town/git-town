package hosting

import (
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/giturl"
	"github.com/git-town/git-town/v13/src/hosting/bitbucket"
	"github.com/git-town/git-town/v13/src/hosting/gitea"
	"github.com/git-town/git-town/v13/src/hosting/github"
	"github.com/git-town/git-town/v13/src/hosting/gitlab"
)

func Detect(originURL *giturl.Parts, hostingPlatform configdomain.HostingPlatform) configdomain.HostingPlatform {
	switch {
	case bitbucket.Detect(originURL, hostingPlatform):
		return configdomain.HostingPlatformBitbucket
	case gitea.Detect(originURL, hostingPlatform):
		return configdomain.HostingPlatformGitea
	case github.Detect(originURL, hostingPlatform):
		return configdomain.HostingPlatformGitHub
	case gitlab.Detect(originURL, hostingPlatform):
		return configdomain.HostingPlatformGitLab
	default:
		return configdomain.HostingPlatformNone
	}
}
