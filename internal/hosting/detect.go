package hosting

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/hosting/bitbucketcloud"
	"github.com/git-town/git-town/v16/internal/hosting/bitbucketdatacenter"
	"github.com/git-town/git-town/v16/internal/hosting/gitea"
	"github.com/git-town/git-town/v16/internal/hosting/github"
	"github.com/git-town/git-town/v16/internal/hosting/gitlab"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

func Detect(remoteURL giturl.Parts, userOverride Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userOverride.IsSome() {
		return userOverride
	}
	detectors := map[configdomain.HostingPlatform]func(giturl.Parts) bool{
		configdomain.HostingPlatformBitbucket:           bitbucketcloud.Detect,
		configdomain.HostingPlatformBitbucketDatacenter: bitbucketdatacenter.Detect,
		configdomain.HostingPlatformGitea:               gitea.Detect,
		configdomain.HostingPlatformGitHub:              github.Detect,
		configdomain.HostingPlatformGitLab:              gitlab.Detect,
	}
	for platform, detector := range detectors {
		if detector(remoteURL) {
			return Some(platform)
		}
	}
	return None[configdomain.HostingPlatform]()
}
