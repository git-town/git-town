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

func Load(env EnvVars) (configdomain.PartialConfig, error) {
	autoResolve, errAutoResolve := gohacks.ParseBoolOpt[configdomain.AutoResolve](env.Get(autoResolve), autoResolve)
	autoSync, errAutoSync := gohacks.ParseBoolOpt[configdomain.AutoSync](env.Get(autoSync), autoSync)
	contributionRegex, errContribRegex := configdomain.ParseContributionRegex(env.Get(contributionRegex))
	detached, errDetached := gohacks.ParseBoolOpt[configdomain.Detached](env.Get(detached), detached)
	displayTypesOpt, errDisplayTypes := configdomain.ParseDisplayTypes(env.Get(displayTypes), displayTypes)
	dryRun, errDryRun := gohacks.ParseBoolOpt[configdomain.DryRun](env.Get(dryRun), dryRun)
	featureRegex, errFeatureRegex := configdomain.ParseFeatureRegex(env.Get(featureRegex))
	forgeType, errForgeType := forgedomain.ParseForgeType(env.Get(forgeType))
	gitAuthorEmailValue := NewOption(gitdomain.GitUserEmail(env.Get(gitAuthorEmail)))
	gitCommitterEmailValue := NewOption(gitdomain.GitUserEmail(env.Get(gitCommitterEmail)))
	gitUserEmail := gitAuthorEmailValue.Or(gitCommitterEmailValue)
	gitAuthorNameValue := NewOption(gitdomain.GitUserName(env.Get(gitAuthorName)))
	gitCommitterNameValue := NewOption(gitdomain.GitUserName(env.Get(gitCommitterName)))
	gitUserName := gitAuthorNameValue.Or(gitCommitterNameValue)
	githubConnectorType, errGitHubConnectorType := forgedomain.ParseGitHubConnectorType(env.Get(githubConnectorType))
	gitlabConnectorType, errGitLabConnectorType := forgedomain.ParseGitLabConnectorType(env.Get(gitlabConnectorType))
	newBranchType, errNewBranchType := configdomain.ParseBranchType(env.Get(newBranchType))
	observedRegex, errObservedRegex := configdomain.ParseObservedRegex(env.Get(observedRegex))
	order, errOrder := configdomain.ParseOrder(env.Get(order), order)
	offline, errOffline := gohacks.ParseBoolOpt[configdomain.Offline](env.Get(offline), offline)
	perennialRegex, errPerennialRegex := configdomain.ParsePerennialRegex(env.Get(perennialRegex))
	proposalsShowLineage, errProposalsShowLineage := forgedomain.ParseProposalsShowLineage(env.Get(proposalsShowLineage))
	pushBranches, errPushBranches := gohacks.ParseBoolOpt[configdomain.PushBranches](env.Get(pushBranches), pushBranches)
	pushHook, errPushHook := gohacks.ParseBoolOpt[configdomain.PushHook](env.Get(pushHook), pushHook)
	shareNewBranches, errShareNewBranches := configdomain.ParseShareNewBranches(env.Get(shareNewBranches), shareNewBranches)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := gohacks.ParseBoolOpt[configdomain.ShipDeleteTrackingBranch](env.Get(shipDeleteTrackingBranch), shipDeleteTrackingBranch)
	shipStrategy, errShipStrategy := configdomain.ParseShipStrategy(env.Get(shipStrategy))
	stash, errStash := gohacks.ParseBoolOpt[configdomain.Stash](env.Get(stash), stash)
	syncFeatureStrategy, errSyncFeatureStrategy := configdomain.ParseSyncFeatureStrategy(env.Get(syncFeatureStrategy))
	syncPerennialStrategy, errSyncPerennialStrategy := configdomain.ParseSyncPerennialStrategy(env.Get(syncPerennialStrategy))
	syncPrototypeStrategy, errSyncPrototypeStrategy := configdomain.ParseSyncPrototypeStrategy(env.Get(syncPrototypeStrategy))
	syncTags, errSyncTags := gohacks.ParseBoolOpt[configdomain.SyncTags](env.Get(syncTags), syncTags)
	syncUpstream, errSyncUpstream := gohacks.ParseBoolOpt[configdomain.SyncUpstream](env.Get(syncUpstream), syncUpstream)
	unknownBranchType, errUnknownBranchType := configdomain.ParseBranchType(env.Get(unknownBranchType))
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
		GitHubUsername:           None[forgedomain.GitHubUsername](), // GitHub username is not loaded from env vars
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
