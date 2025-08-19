package envconfig

import (
	"cmp"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	autoResolve      = "GITHUB_AUTO_RESOLVE"
	dryRun           = "GIT_TOWN_DRY_RUN"
	offline          = "GIT_TOWN_OFFLINE"
	pushHook         = "GIT_TOWN_PUSH_HOOK"
	shareNewBranches = "GIT_TOWN_SHARE_NEW_BRANCHES"
)

func Load(env Environment) (configdomain.PartialConfig, error) {
	autoResolve, errAutoResolve := configdomain.ParseAutoResolve(env.Get(autoResolve), autoResolve)
	contributionRegex, errContribRegex := configdomain.ParseContributionRegex(env.Get("GIT_TOWN_CONTRIBUTION_REGEX"))
	dryRun, errDryRun := configdomain.ParseDryRun(env.Get(dryRun), dryRun)
	featureRegex, errFeatureRegex := configdomain.ParseFeatureRegex(env.Get("GIT_TOWN_FEATURE_REGEX"))
	forgeType, errForgeType := forgedomain.ParseForgeType(env.Get("GIT_TOWN_FORGE_TYPE"))
	githubConnectorType, errGitHubConnectorType := forgedomain.ParseGitHubConnectorType(env.Get("GIT_TOWN_GITHUB_CONNECTOR_TYPE"))
	gitlabConnectorType, errGitLabConnectorType := forgedomain.ParseGitLabConnectorType(env.Get("GIT_TOWN_GITLAB_CONNECTOR_TYPE"))
	newBranchType, errNewBranchType := configdomain.ParseBranchType(env.Get("GIT_TOWN_NEW_BRANCH_TYPE"))
	observedRegex, errObservedRegex := configdomain.ParseObservedRegex(env.Get("GIT_TOWN_OBSERVED_REGEX"))
	offline, errOffline := configdomain.ParseOffline(env.Get(offline), offline)
	perennialRegex, errPerennialRegex := configdomain.ParsePerennialRegex(env.Get("GIT_TOWN_PERENNIAL_REGEX"))
	proposalsShowLineage, errProposalsShowLineage := forgedomain.ParseProposalsShowLineage(env.Get("GIT_TOWN_PROPOSALS_SHOW_LINEAGE"))
	pushHook, errPushHook := configdomain.ParsePushHook(env.Get(pushHook), pushHook)
	shareNewBranches, errShareNewBranches := configdomain.ParseShareNewBranches(env.Get(shareNewBranches), shareNewBranches)
	err := cmp.Or(
		errAutoResolve,
		errContribRegex,
		errDryRun,
		errFeatureRegex,
		errForgeType,
		errGitHubConnectorType,
		errGitLabConnectorType,
		errNewBranchType,
		errObservedRegex,
		errOffline,
		errPerennialRegex,
		errProposalsShowLineage,
		errPushHook,
		errShareNewBranches,
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
		GitHubConnectorType:      githubConnectorType,
		GitHubToken:              forgedomain.ParseGitHubToken(env.Get("GIT_TOWN_GITHUB_TOKEN", "GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")),
		GitLabConnectorType:      gitlabConnectorType,
		GitLabToken:              forgedomain.ParseGitLabToken(env.Get("GIT_TOWN_GITLAB_TOKEN")),
		GitUserEmail:             None[gitdomain.GitUserEmail](), // not loaded from env vars
		GitUserName:              None[gitdomain.GitUserName](),  // not loaded from env vars
		GiteaToken:               forgedomain.ParseGiteaToken(env.Get("GIT_TOWN_GITEA_TOKEN")),
		HostingOriginHostname:    configdomain.ParseHostingOriginHostname(env.Get("GIT_TOWN_ORIGIN_HOSTNAME")),
		Lineage:                  configdomain.NewLineage(), // not loaded from env vars
		MainBranch:               gitdomain.NewLocalBranchNameOption(env.Get("GIT_TOWN_MAIN_BRANCH")),
		NewBranchType:            configdomain.NewBranchTypeOpt(newBranchType),
		ObservedRegex:            observedRegex,
		Offline:                  offline,
		PerennialBranches:        gitdomain.ParseLocalBranchNames(env.Get("GIT_TOWN_PERENNIAL_BRANCHES")),
		PerennialRegex:           perennialRegex,
		ProposalsShowLineage:     proposalsShowLineage,
		PushHook:                 pushHook,
		ShareNewBranches:         shareNewBranches,
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
