package envconfig

import (
	"cmp"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	autoResolve = "GITHUB_AUTO_RESOLVE"
	dryRun      = "GIT_TOWN_DRY_RUN"
)

func Load(env Environment) (configdomain.PartialConfig, error) {
	autoResolve, errAutoResolve := configdomain.ParseAutoResolve(env.Get(autoResolve), autoResolve)
	contributionRegex, errContribRegex := configdomain.ParseContributionRegex(env.Get("GIT_TOWN_CONTRIBUTION_REGEX"))
	dryRun, errDryRun := configdomain.ParseDryRun(env.Get(dryRun), dryRun)
	featureRegex, errFeatureRegex := configdomain.ParseFeatureRegex(env.Get("GIT_TOWN_FEATURE_REGEX"))
	forgeType, errForgeType := forgedomain.ParseForgeType(env.Get("GIT_TOWN_FORGE_TYPE"))
	err := cmp.Or(
		errAutoResolve,
		errContribRegex,
		errDryRun,
		errFeatureRegex,
		errForgeType,
	)
	return configdomain.PartialConfig{
		Aliases:                  configdomain.Aliases{}, // aliases aren't loaded from env vars
		AutoResolve:              autoResolve,
		BitbucketAppPassword:     forgedomain.ParseBitbucketAppPassword(env.Get("GIT_TOWN_BITBUCKET_APP_PASSWORD")),
		BitbucketUsername:        forgedomain.ParseBitbucketUsername(env.Get("GIT_TOWN_BITBUCKET_USERNAME")),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{}, // not loaded from env vars
		CodebergToken:            forgedomain.ParseCodebergToken(env.Get("GIT_TOWN_CODEBERG_TOKEN")),
		ContributionRegex:        contributionRegex,
		DevRemote:                gitdomain.NewRemote(env.Get("GIT_TOWN_DEV_REMOTE")),
		DryRun:                   dryRun,
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
		GitHubToken:              forgedomain.ParseGitHubToken(env.Get("GIT_TOWN_GITHUB_TOKEN", "GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")),
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
		SyncFeatureStrategy:      None[configdomain.SyncFeatureStrategy](),
		SyncPerennialStrategy:    None[configdomain.SyncPerennialStrategy](),
		SyncPrototypeStrategy:    None[configdomain.SyncPrototypeStrategy](),
		SyncTags:                 None[configdomain.SyncTags](),
		SyncUpstream:             None[configdomain.SyncUpstream](),
		UnknownBranchType:        None[configdomain.UnknownBranchType](),
		Verbose:                  None[configdomain.Verbose](),
	}, err
}
