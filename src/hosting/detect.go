package hosting

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/bitbucket"
	"github.com/git-town/git-town/v11/src/hosting/gitea"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/hosting/gitlab"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
)

func detect(originURL *giturl.Parts, hostingPlatform configdomain.HostingPlatform) hostingdomain.Platform {
	switch {
	case bitbucket.Detect(originURL, hostingPlatform):
		return hostingdomain.PlatformBitbucket
	case gitea.Detect(originURL, hostingPlatform):
		return hostingdomain.PlatformGitea
	case github.Detect(originURL, hostingPlatform):
		return hostingdomain.PlatformGithub
	case gitlab.Detect(originURL, hostingPlatform):
		return hostingdomain.PlatformGitlab
	default:
		return hostingdomain.PlatformNone
	}
}
