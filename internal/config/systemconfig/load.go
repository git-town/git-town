package systemconfig

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func Load() configdomain.PartialConfig {
	hasTTY := DetermineTTY()
	displayDialogs := configdomain.NewDisplayDialogsFromTTY(hasTTY)
	return configdomain.PartialConfig{
		Aliases:                     configdomain.Aliases{},
		AutoResolve:                 None[configdomain.AutoResolve](),
		AutoSync:                    None[configdomain.AutoSync](),
		BitbucketAppPassword:        None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:           None[forgedomain.BitbucketUsername](),
		BranchPrefix:                None[configdomain.BranchPrefix](),
		BranchTypeOverrides:         configdomain.BranchTypeOverrides{},
		Browser:                     None[configdomain.Browser](),
		ContributionRegex:           None[configdomain.ContributionRegex](),
		Detached:                    None[configdomain.Detached](),
		DevRemote:                   None[gitdomain.Remote](),
		DisplayDialogs:              displayDialogs,
		DisplayTypes:                None[configdomain.DisplayTypes](),
		DryRun:                      None[configdomain.DryRun](),
		FeatureRegex:                None[configdomain.FeatureRegex](),
		ForgeType:                   None[forgedomain.ForgeType](),
		ForgejoToken:                None[forgedomain.ForgejoToken](),
		GitUserEmail:                None[gitdomain.GitUserEmail](),
		GitUserName:                 None[gitdomain.GitUserName](),
		GiteaToken:                  None[forgedomain.GiteaToken](),
		GithubConnectorType:         None[forgedomain.GithubConnectorType](),
		GithubToken:                 None[forgedomain.GithubToken](),
		GitlabConnectorType:         None[forgedomain.GitlabConnectorType](),
		GitlabToken:                 None[forgedomain.GitlabToken](),
		HostingOriginHostname:       None[configdomain.HostingOriginHostname](),
		IgnoreUncommitted:           None[configdomain.IgnoreUncommitted](),
		Lineage:                     configdomain.NewLineage(),
		MainBranch:                  None[gitdomain.LocalBranchName](),
		NewBranchType:               None[configdomain.NewBranchType](),
		ObservedRegex:               None[configdomain.ObservedRegex](),
		Offline:                     None[configdomain.Offline](),
		Order:                       None[configdomain.Order](),
		PerennialBranches:           gitdomain.LocalBranchNames{},
		PerennialRegex:              None[configdomain.PerennialRegex](),
		ProposalBreadcrumb:          None[configdomain.ProposalBreadcrumb](),
		ProposalBreadcrumbDirection: None[configdomain.ProposalBreadcrumbDirection](),
		ProposeHeadless:             None[forgedomain.ProposeHeadless](),
		PushBranches:                None[configdomain.PushBranches](),
		PushHook:                    None[configdomain.PushHook](),
		ShareNewBranches:            None[configdomain.ShareNewBranches](),
		ShipDeleteTrackingBranch:    None[configdomain.ShipDeleteTrackingBranch](),
		ShipStrategy:                None[configdomain.ShipStrategy](),
		Stash:                       None[configdomain.Stash](),
		SyncFeatureStrategy:         None[configdomain.SyncFeatureStrategy](),
		SyncPerennialStrategy:       None[configdomain.SyncPerennialStrategy](),
		SyncPrototypeStrategy:       None[configdomain.SyncPrototypeStrategy](),
		SyncTags:                    None[configdomain.SyncTags](),
		SyncUpstream:                None[configdomain.SyncUpstream](),
		UnknownBranchType:           None[configdomain.UnknownBranchType](),
		Verbose:                     None[configdomain.Verbose](),
	}
}
