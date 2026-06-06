package forge

import (
	"github.com/git-town/git-town/v23/internal/forge/azuredevops"
	"github.com/git-town/git-town/v23/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v23/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/forge/forgejo"
	"github.com/git-town/git-town/v23/internal/forge/gitea"
	"github.com/git-town/git-town/v23/internal/forge/github"
	"github.com/git-town/git-town/v23/internal/forge/gitlab"
	"github.com/git-town/git-town/v23/internal/git/giturl"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

func Detect(remoteURL giturl.Parts, userOverride Option[forgedomain.ForgeType]) Option[forgedomain.DetectedForgeType] {
	if override, hasOverride := userOverride.Get(); hasOverride {
		return Some(override.Detected())
	}
	detectors := []detector{
		{forgedomain.ForgeTypeAzuredevops, azuredevops.Detect},
		{forgedomain.ForgeTypeBitbucket, bitbucketcloud.Detect},
		{forgedomain.ForgeTypeBitbucketDatacenter, bitbucketdatacenter.Detect},
		{forgedomain.ForgeTypeForgejo, forgejo.Detect},
		{forgedomain.ForgeTypeGitea, gitea.Detect},
		{forgedomain.ForgeTypeGithub, github.Detect},
		{forgedomain.ForgeTypeGitlab, gitlab.Detect},
	}
	for _, detector := range detectors {
		if detector.implementation(remoteURL) {
			return Some(forgedomain.DetectedForgeType(detector.forgeType))
		}
	}
	return None[forgedomain.DetectedForgeType]()
}

type detector struct {
	forgeType      forgedomain.ForgeType
	implementation func(giturl.Parts) bool
}
