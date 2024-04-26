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
	detectors := map[configdomain.HostingPlatform]func(*giturl.Parts) bool{
		configdomain.HostingPlatformBitbucket: bitbucket.Detect,
		configdomain.HostingPlatformGitea:     gitea.Detect,
		configdomain.HostingPlatformGitHub:    github.Detect,
		configdomain.HostingPlatformGitLab:    gitlab.Detect,
	}
	for platform, detector := range detectors {
		if detector(originURL) {
			return Some(platform)
		}
	}
	return None[configdomain.HostingPlatform]()
}
