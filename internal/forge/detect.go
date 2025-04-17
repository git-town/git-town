package forge

import (
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v19/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v19/internal/forge/codeberg"
	"github.com/git-town/git-town/v19/internal/forge/gitea"
	"github.com/git-town/git-town/v19/internal/forge/github"
	"github.com/git-town/git-town/v19/internal/forge/gitlab"
	"github.com/git-town/git-town/v19/internal/git/giturl"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

func Detect(remoteURL giturl.Parts, userOverride Option[configdomain.ForgeType]) Option[configdomain.ForgeType] {
	if userOverride.IsSome() {
		return userOverride
	}
	detectors := map[configdomain.ForgeType]func(giturl.Parts) bool{
		configdomain.ForgeTypeBitbucket:           bitbucketcloud.Detect,
		configdomain.ForgeTypeBitbucketDatacenter: bitbucketdatacenter.Detect,
		configdomain.ForgeTypeCodeberg:            codeberg.Detect,
		configdomain.ForgeTypeGitea:               gitea.Detect,
		configdomain.ForgeTypeGitHub:              github.Detect,
		configdomain.ForgeTypeGitLab:              gitlab.Detect,
	}
	for platform, detector := range detectors {
		if detector(remoteURL) {
			return Some(platform)
		}
	}
	return None[configdomain.ForgeType]()
}
