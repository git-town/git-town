package envconfig

import (
	"cmp"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	autoResolve              = "GITHUB_AUTO_RESOLVE"
	contributionRegex        = "GIT_TOWN_CONTRIBUTION_REGEX"
	dryRun                   = "GIT_TOWN_DRY_RUN"
	featureRegex             = "GIT_TOWN_FEATURE_REGEX"
	forgeType                = "GIT_TOWN_FORGE_TYPE"
	githubConnectorType      = "GIT_TOWN_GITHUB_CONNECTOR_TYPE"
	gitlabConnectorType      = "GIT_TOWN_GITLAB_CONNECTOR_TYPE"
	newBranchType            = "GIT_TOWN_NEW_BRANCH_TYPE"
	observedRegex            = "GIT_TOWN_OBSERVED_REGEX"
	offline                  = "GIT_TOWN_OFFLINE"
	perennialRegex           = "GIT_TOWN_PERENNIAL_REGEX"
	proposalsShowLineage     = "GIT_TOWN_PROPOSALS_SHOW_LINEAGE"
	pushHook                 = "GIT_TOWN_PUSH_HOOK"
	shareNewBranches         = "GIT_TOWN_SHARE_NEW_BRANCHES"
	shipDeleteTrackingBranch = "GIT_TOWN_SHIP_DELETE_TRACKING_BRANCH"
	shipStrategy             = "GIT_TOWN_SHIP_STRATEGY"
	syncTags                 = "GIT_TOWN_SYNC_TAGS"
	syncUpstream             = "GIT_TOWN_SYNC_UPSTREAM"
	verbose                  = "GIT_TOWN_VERBOSE"
)

func Load(env Environment) (configdomain.PartialConfig, error) {
	autoResolve, errAutoResolve := gohacks.ParseBoolOpt[configdomain.AutoResolve](env.Get(autoResolve), autoResolve)
	contributionRegex, errContribRegex := configdomain.ParseContributionRegex(env.Get(contributionRegex))
	dryRun, errDryRun := gohacks.ParseBoolOpt[configdomain.DryRun](env.Get(dryRun), dryRun)
	featureRegex, errFeatureRegex := configdomain.ParseFeatureRegex(env.Get(featureRegex))
	forgeType, errForgeType := forgedomain.ParseForgeType(env.Get(forgeType))
	githubConnectorType, errGitHubConnectorType := forgedomain.ParseGitHubConnectorType(env.Get(githubConnectorType))
	gitlabConnectorType, errGitLabConnectorType := forgedomain.ParseGitLabConnectorType(env.Get(gitlabConnectorType))
	newBranchType, errNewBranchType := configdomain.ParseBranchType(env.Get(newBranchType))
	observedRegex, errObservedRegex := configdomain.ParseObservedRegex(env.Get(observedRegex))
	offline, errOffline := gohacks.ParseBoolOpt[configdomain.Offline](env.Get(offline), offline)
	perennialRegex, errPerennialRegex := configdomain.ParsePerennialRegex(env.Get(perennialRegex))
	proposalsShowLineage, errProposalsShowLineage := forgedomain.ParseProposalsShowLineage(env.Get(proposalsShowLineage))
	pushHook, errPushHook := gohacks.ParseBoolOpt[configdomain.PushHook](env.Get(pushHook), pushHook)
	shareNewBranches, errShareNewBranches := configdomain.ParseShareNewBranches(env.Get(shareNewBranches), shareNewBranches)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := gohacks.ParseBoolOpt[configdomain.ShipDeleteTrackingBranch](env.Get(shipDeleteTrackingBranch), shipDeleteTrackingBranch)
	shipStrategy, errShipStrategy := configdomain.ParseShipStrategy(env.Get(shipStrategy))
	syncFeatureStrategy, errSyncFeatureStrategy := configdomain.ParseSyncFeatureStrategy(env.Get("GIT_TOWN_SYNC_FEATURE_STRATEGY"))
	syncPerennialStrategy, errSyncPerennialStrategy := configdomain.ParseSyncPerennialStrategy(env.Get("GIT_TOWN_SYNC_PERENNIAL_STRATEGY"))
	syncPrototypeStrategy, errSyncPrototypeStrategy := configdomain.ParseSyncPrototypeStrategy(env.Get("GIT_TOWN_SYNC_PROTOTYPE_STRATEGY"))
	syncTags, errSyncTags := gohacks.ParseBoolOpt[configdomain.SyncTags](env.Get(syncTags), syncTags)
	syncUpstream, errSyncUpstream := gohacks.ParseBoolOpt[configdomain.SyncUpstream](env.Get(syncUpstream), syncUpstream)
	unknownBranchType, errUnknownBranchType := configdomain.ParseBranchType(env.Get("GIT_TOWN_UNKNOWN_BRANCH_TYPE"))
	verbose, errVerbose := gohacks.ParseBoolOpt[configdomain.Verbose](env.Get(verbose), verbose)
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
		errShipDeleteTrackingBranch,
		errShipStrategy,
		errSyncFeatureStrategy,
		errSyncPerennialStrategy,
		errSyncPrototypeStrategy,
		errSyncTags,
		errSyncUpstream,
		errUnknownBranchType,
		errVerbose,
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
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
		UnknownBranchType:        configdomain.UnknownBranchTypeOpt(unknownBranchType),
		Verbose:                  verbose,
	}, err
}
