package cliconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type NewArgs struct {
	AutoResolve  Option[configdomain.AutoResolve]
	AutoSync     Option[configdomain.AutoSync]
	Detached     Option[configdomain.Detached]
	DryRun       Option[configdomain.DryRun]
	PushBranches Option[configdomain.PushBranches]
	Stash        Option[configdomain.Stash]
	Verbose      Option[configdomain.Verbose]
}

func New(args NewArgs) configdomain.PartialConfig {
	return configdomain.PartialConfig{
		Aliases:                  configdomain.Aliases{},
		AutoResolve:              args.AutoResolve,
		AutoSync:                 args.AutoSync,
		BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:        None[forgedomain.BitbucketUsername](),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
		ForgejoToken:             None[forgedomain.ForgejoToken](),
		ContributionRegex:        None[configdomain.ContributionRegex](),
		Detached:                 args.Detached,
		DevRemote:                None[gitdomain.Remote](),
		DryRun:                   args.DryRun,
		FeatureRegex:             None[configdomain.FeatureRegex](),
		ForgeType:                None[forgedomain.ForgeType](),
		GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
		GitHubToken:              None[forgedomain.GitHubToken](),
		GitLabConnectorType:      None[forgedomain.GitLabConnectorType](),
		GitLabToken:              None[forgedomain.GitLabToken](),
		GitUserEmail:             None[gitdomain.GitUserEmail](),
		GitUserName:              None[gitdomain.GitUserName](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
		Lineage:                  configdomain.NewLineage(),
		MainBranch:               None[gitdomain.LocalBranchName](),
		NewBranchType:            None[configdomain.NewBranchType](),
		ObservedRegex:            None[configdomain.ObservedRegex](),
		Offline:                  None[configdomain.Offline](),
		PerennialBranches:        gitdomain.LocalBranchNames{},
		PerennialRegex:           None[configdomain.PerennialRegex](),
		ProposalsShowLineage:     None[forgedomain.ProposalsShowLineage](),
		PushHook:                 None[configdomain.PushHook](),
		ShareNewBranches:         None[configdomain.ShareNewBranches](),
		ShipDeleteTrackingBranch: None[configdomain.ShipDeleteTrackingBranch](),
		ShipStrategy:             None[configdomain.ShipStrategy](),
		Stash:                    args.Stash,
		SyncFeatureStrategy:      None[configdomain.SyncFeatureStrategy](),
		SyncPerennialStrategy:    None[configdomain.SyncPerennialStrategy](),
		SyncPrototypeStrategy:    None[configdomain.SyncPrototypeStrategy](),
		PushBranches:             args.PushBranches,
		SyncTags:                 None[configdomain.SyncTags](),
		SyncUpstream:             None[configdomain.SyncUpstream](),
		UnknownBranchType:        None[configdomain.UnknownBranchType](),
		Verbose:                  args.Verbose,
	}
}
