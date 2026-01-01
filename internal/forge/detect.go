package forge

import (
	"github.com/git-town/git-town/v22/internal/forge/azuredevops"
	"github.com/git-town/git-town/v22/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v22/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/forgejo"
	"github.com/git-town/git-town/v22/internal/forge/gitea"
	"github.com/git-town/git-town/v22/internal/forge/github"
	"github.com/git-town/git-town/v22/internal/forge/gitlab"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func Detect(remoteURL giturl.Parts, userOverride Option[forgedomain.ForgeType]) Option[forgedomain.ForgeType] {
	if userOverride.IsSome() {
		return userOverride
	}
	detectors := []detector{
		{forgedomain.ForgeTypeAzureDevOps, azuredevops.Detect},
		{forgedomain.ForgeTypeBitbucket, bitbucketcloud.Detect},
		{forgedomain.ForgeTypeBitbucketDatacenter, bitbucketdatacenter.Detect},
		{forgedomain.ForgeTypeForgejo, forgejo.Detect},
		{forgedomain.ForgeTypeGitea, gitea.Detect},
		{forgedomain.ForgeTypeGithub, github.Detect},
		{forgedomain.ForgeTypeGitLab, gitlab.Detect},
	}
	for _, detector := range detectors {
		if detector.implementation(remoteURL) {
			return Some(detector.forgeType)
		}
	}
	return None[forgedomain.ForgeType]()
}

type detector struct {
	forgeType      forgedomain.ForgeType
	implementation func(giturl.Parts) bool
}
