package hosting

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/hosting/bitbucket"
	"github.com/git-town/git-town/v14/src/hosting/gitea"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/git-town/git-town/v14/src/hosting/gitlab"
)

func Detect(originURL *giturl.Parts, hostingPlatform gohacks.Option[configdomain.HostingPlatform]) gohacks.Option[configdomain.HostingPlatform] {
	switch {
	case bitbucket.Detect(originURL, hostingPlatform):
		return gohacks.NewOption(configdomain.HostingPlatformBitbucket)
	case gitea.Detect(originURL, hostingPlatform):
		return gohacks.NewOption(configdomain.HostingPlatformGitea)
	case github.Detect(originURL, hostingPlatform):
		return gohacks.NewOption(configdomain.HostingPlatformGitHub)
	case gitlab.Detect(originURL, hostingPlatform):
		return gohacks.NewOption(configdomain.HostingPlatformGitLab)
	default:
		return gohacks.NewOptionNone[configdomain.HostingPlatform]()
	}
}
