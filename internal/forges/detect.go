package forges

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/forges/bitbucketcloud"
	"github.com/git-town/git-town/v18/internal/forges/bitbucketdatacenter"
	"github.com/git-town/git-town/v18/internal/forges/gitea"
	"github.com/git-town/git-town/v18/internal/forges/github"
	"github.com/git-town/git-town/v18/internal/forges/gitlab"
	"github.com/git-town/git-town/v18/internal/git/giturl"
	. "github.com/git-town/git-town/v18/pkg/prelude"
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
