package hosting

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting/bitbucket"
	"github.com/git-town/git-town/v14/src/hosting/gitea"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/git-town/git-town/v14/src/hosting/gitlab"
)

func Detect(originURL *giturl.Parts, userOverride Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userOverride.IsSome() {
		return userOverride
	}
	switch {
	case bitbucket.Detect(originURL):
		return Some(configdomain.HostingPlatformBitbucket)
	case gitea.Detect(originURL):
		return Some(configdomain.HostingPlatformGitea)
	case github.Detect(originURL):
		return Some(configdomain.HostingPlatformGitHub)
	case gitlab.Detect(originURL):
		return Some(configdomain.HostingPlatformGitLab)
	default:
		return None[configdomain.HostingPlatform]()
	}
}
