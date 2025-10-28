package envconfig

import (
	"cmp"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	autoResolve              = "GIT_TOWN_AUTO_RESOLVE"
	autoSync                 = "GIT_TOWN_AUTO_SYNC"
	bitbucketAppPassword     = "GIT_TOWN_BITBUCKET_APP_PASSWORD"
	bitbucketUserName        = "GIT_TOWN_BITBUCKET_USERNAME"
	forgejoToken             = "GIT_TOWN_FORGEJO_TOKEN"
	contributionRegex        = "GIT_TOWN_CONTRIBUTION_REGEX"
	detached                 = "GIT_TOWN_DETACHED"
	devRemote                = "GIT_TOWN_DEV_REMOTE"
	displayTypes             = "GIT_TOWN_DISPLAY_TYPES"
	dryRun                   = "GIT_TOWN_DRY_RUN"
	featureRegex             = "GIT_TOWN_FEATURE_REGEX"
	forgeType                = "GIT_TOWN_FORGE_TYPE"
	giteaToken               = "GIT_TOWN_GITEA_TOKEN"
	gitAuthorEmail           = "GIT_AUTHOR_EMAIL"
	gitAuthorName            = "GIT_AUTHOR_NAME"
	gitCommitterEmail        = "GIT_COMMITTER_EMAIL"
	gitCommitterName         = "GIT_COMMITTER_NAME"
	githubConnectorType      = "GIT_TOWN_GITHUB_CONNECTOR_TYPE"
	githubToken              = "GIT_TOWN_GITHUB_TOKEN"
	gitlabConnectorType      = "GIT_TOWN_GITLAB_CONNECTOR_TYPE"
	gitlabToken              = "GIT_TOWN_GITLAB_TOKEN"
	mainBranch               = "GIT_TOWN_MAIN_BRANCH"
	newBranchType            = "GIT_TOWN_NEW_BRANCH_TYPE"
	observedRegex            = "GIT_TOWN_OBSERVED_REGEX"
	order                    = "GIT_TOWN_ORDER"
	originHostname           = "GIT_TOWN_ORIGIN_HOSTNAME"
	offline                  = "GIT_TOWN_OFFLINE"
	perennialBranches        = "GIT_TOWN_PERENNIAL_BRANCHES"
	perennialRegex           = "GIT_TOWN_PERENNIAL_REGEX"
	proposalsShowLineage     = "GIT_TOWN_PROPOSALS_SHOW_LINEAGE"
	pushBranches             = "GIT_TOWN_PUSH_BRANCHES"
	pushHook                 = "GIT_TOWN_PUSH_HOOK"
	shareNewBranches         = "GIT_TOWN_SHARE_NEW_BRANCHES"
	shipDeleteTrackingBranch = "GIT_TOWN_SHIP_DELETE_TRACKING_BRANCH"
	shipStrategy             = "GIT_TOWN_SHIP_STRATEGY"
	stash                    = "GIT_TOWN_STASH"
	syncFeatureStrategy      = "GIT_TOWN_SYNC_FEATURE_STRATEGY"
	syncPerennialStrategy    = "GIT_TOWN_SYNC_PERENNIAL_STRATEGY"
	syncPrototypeStrategy    = "GIT_TOWN_SYNC_PROTOTYPE_STRATEGY"
	syncTags                 = "GIT_TOWN_SYNC_TAGS"
	syncUpstream             = "GIT_TOWN_SYNC_UPSTREAM"
	unknownBranchType        = "GIT_TOWN_UNKNOWN_BRANCH_TYPE"
	verbose                  = "GIT_TOWN_VERBOSE"
)

func load[T any](env EnvVars, varName string, parser func(string, string) (T, error)) (T, error) {
	return parser(env.Get(varName), varName)
}

func Load(env EnvVars) (configdomain.PartialConfig, error) {
	autoResolve, errAutoResolve := load(env, autoResolve, gohacks.ParseBoolOpt[configdomain.AutoResolve])
	autoSync, errAutoSync := load(env, autoSync, gohacks.ParseBoolOpt[configdomain.AutoSync])
	contributionRegex, errContribRegex := load(env, contributionRegex, configdomain.ParseContributionRegex)
	detached, errDetached := load(env, detached, gohacks.ParseBoolOpt[configdomain.Detached])
	displayTypesOpt, errDisplayTypes := load(env, displayTypes, configdomain.ParseDisplayTypes)
	dryRun, errDryRun := load(env, dryRun, gohacks.ParseBoolOpt[configdomain.DryRun])
	featureRegex, errFeatureRegex := load(env, featureRegex, configdomain.ParseFeatureRegex)
	forgeType, errForgeType := load(env, forgeType, forgedomain.ParseForgeType)
	gitAuthorEmailValue := NewOption(gitdomain.GitUserEmail(env.Get(gitAuthorEmail)))
	gitCommitterEmailValue := NewOption(gitdomain.GitUserEmail(env.Get(gitCommitterEmail)))
	gitUserEmail := gitAuthorEmailValue.Or(gitCommitterEmailValue)
	gitAuthorNameValue := NewOption(gitdomain.GitUserName(env.Get(gitAuthorName)))
	gitCommitterNameValue := NewOption(gitdomain.GitUserName(env.Get(gitCommitterName)))
	gitUserName := gitAuthorNameValue.Or(gitCommitterNameValue)
	githubConnectorType, errGitHubConnectorType := load(env, githubConnectorType, forgedomain.ParseGitHubConnectorType)
	gitlabConnectorType, errGitLabConnectorType := load(env, gitlabConnectorType, forgedomain.ParseGitLabConnectorType)
	newBranchType, errNewBranchType := load(env, newBranchType, configdomain.ParseBranchType)
	observedRegex, errObservedRegex := load(env, observedRegex, configdomain.ParseObservedRegex)
	order, errOrder := configdomain.ParseOrder(env.Get(order), order)
	offline, errOffline := load(env, offline, gohacks.ParseBoolOpt[configdomain.Offline])
	perennialRegex, errPerennialRegex := load(env, perennialRegex, configdomain.ParsePerennialRegex)
	proposalsShowLineage, errProposalsShowLineage := load(env, proposalsShowLineage, forgedomain.ParseProposalsShowLineage)
	pushBranches, errPushBranches := load(env, pushBranches, gohacks.ParseBoolOpt[configdomain.PushBranches])
	pushHook, errPushHook := load(env, pushHook, gohacks.ParseBoolOpt[configdomain.PushHook])
	shareNewBranches, errShareNewBranches := load(env, shareNewBranches, configdomain.ParseShareNewBranches)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := load(env, shipDeleteTrackingBranch, gohacks.ParseBoolOpt[configdomain.ShipDeleteTrackingBranch])
	shipStrategy, errShipStrategy := load(env, shipStrategy, configdomain.ParseShipStrategy)
	stash, errStash := load(env, stash, gohacks.ParseBoolOpt[configdomain.Stash])
	syncFeatureStrategy, errSyncFeatureStrategy := load(env, syncFeatureStrategy, configdomain.ParseSyncFeatureStrategy)
	syncPerennialStrategy, errSyncPerennialStrategy := load(env, syncPerennialStrategy, configdomain.ParseSyncPerennialStrategy)
	syncPrototypeStrategy, errSyncPrototypeStrategy := load(env, syncPrototypeStrategy, configdomain.ParseSyncPrototypeStrategy)
	syncTags, errSyncTags := load(env, syncTags, gohacks.ParseBoolOpt[configdomain.SyncTags])
	syncUpstream, errSyncUpstream := gohacks.ParseBoolOpt[configdomain.SyncUpstream](env.Get(syncUpstream), syncUpstream)
	unknownBranchType, errUnknownBranchType := configdomain.ParseBranchType(env.Get(unknownBranchType), unknownBranchType)
	verbose, errVerbose := gohacks.ParseBoolOpt[configdomain.Verbose](env.Get(verbose), verbose)
	err := cmp.Or(
		errAutoResolve,
		errAutoSync,
		errContribRegex,
		errDetached,
		errDisplayTypes,
		errDryRun,
		errFeatureRegex,
		errForgeType,
		errGitHubConnectorType,
		errGitLabConnectorType,
		errNewBranchType,
		errObservedRegex,
		errOffline,
		errOrder,
		errPerennialRegex,
		errProposalsShowLineage,
		errPushBranches,
		errPushHook,
		errShareNewBranches,
		errShipDeleteTrackingBranch,
		errShipStrategy,
		errStash,
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
		AutoSync:                 autoSync,
		BitbucketAppPassword:     forgedomain.ParseBitbucketAppPassword(env.Get(bitbucketAppPassword)),
		BitbucketUsername:        forgedomain.ParseBitbucketUsername(env.Get(bitbucketUserName)),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{}, // not loaded from env vars
		ForgejoToken:             forgedomain.ParseForgejoToken(env.Get(forgejoToken)),
		ContributionRegex:        contributionRegex,
		Detached:                 detached,
		DevRemote:                gitdomain.NewRemote(env.Get(devRemote)),
		DisplayTypes:             displayTypesOpt,
		DryRun:                   dryRun,
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubConnectorType:      githubConnectorType,
		GitHubToken:              forgedomain.ParseGitHubToken(env.Get(githubToken, "GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")),
		GitLabConnectorType:      gitlabConnectorType,
		GitLabToken:              forgedomain.ParseGitLabToken(env.Get(gitlabToken)),
		GitUserEmail:             gitUserEmail,
		GitUserName:              gitUserName,
		GiteaToken:               forgedomain.ParseGiteaToken(env.Get(giteaToken)),
		HostingOriginHostname:    configdomain.ParseHostingOriginHostname(env.Get(originHostname)),
		Lineage:                  configdomain.NewLineage(), // not loaded from env vars
		MainBranch:               gitdomain.NewLocalBranchNameOption(env.Get(mainBranch)),
		NewBranchType:            configdomain.NewBranchTypeOpt(newBranchType),
		ObservedRegex:            observedRegex,
		Offline:                  offline,
		Order:                    order,
		PerennialBranches:        gitdomain.ParseLocalBranchNames(env.Get(perennialBranches)),
		PerennialRegex:           perennialRegex,
		ProposalsShowLineage:     proposalsShowLineage,
		PushBranches:             pushBranches,
		PushHook:                 pushHook,
		ShareNewBranches:         shareNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		Stash:                    stash,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
		UnknownBranchType:        configdomain.UnknownBranchTypeOpt(unknownBranchType),
		Verbose:                  verbose,
	}, err
}
