package forge

import (
	"github.com/git-town/git-town/v21/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v21/internal/forge/codeberg"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gitea"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func Detect(remoteURL giturl.Parts, userOverride Option[forgedomain.ForgeType]) Option[forgedomain.ForgeType] {
	if userOverride.IsSome() {
		return userOverride
	}
	detectors := map[forgedomain.ForgeType]func(giturl.Parts) bool{
		forgedomain.ForgeTypeBitbucket:           bitbucketcloud.Detect,
		forgedomain.ForgeTypeBitbucketDatacenter: bitbucketdatacenter.Detect,
		forgedomain.ForgeTypeCodeberg:            codeberg.Detect,
		forgedomain.ForgeTypeGitea:               gitea.Detect,
		forgedomain.ForgeTypeGitHub:              github.Detect,
		forgedomain.ForgeTypeGitLab:              gitlab.Detect,
	}
	for forgeType, detector := range detectors {
		if detector(remoteURL) {
			return Some(forgeType)
		}
	}
	return None[forgedomain.ForgeType]()
}
