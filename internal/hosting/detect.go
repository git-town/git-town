package hosting

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting/bitbucket"
	"github.com/git-town/git-town/v15/internal/hosting/gitea"
	"github.com/git-town/git-town/v15/internal/hosting/github"
	"github.com/git-town/git-town/v15/internal/hosting/gitlab"
)

func Detect(remoteURL giturl.Parts, userOverride Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userOverride.IsSome() {
		return userOverride
	}
	detectors := map[configdomain.HostingPlatform]func(giturl.Parts) bool{
		configdomain.HostingPlatformBitbucket: bitbucket.Detect,
		configdomain.HostingPlatformGitea:     gitea.Detect,
		configdomain.HostingPlatformGitHub:    github.Detect,
		configdomain.HostingPlatformGitLab:    gitlab.Detect,
	}
	for platform, detector := range detectors {
		if detector(remoteURL) {
			return Some(platform)
		}
	}
	return None[configdomain.HostingPlatform]()
}
