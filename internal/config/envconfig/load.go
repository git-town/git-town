package envconfig

import (
	"cmp"

	"github.com/git-town/git-town/v23/internal/browser/browserdomain"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
)

const (
	// keep-sorted start
	Browser                     = "BROWSER"
	autoResolve                 = "GIT_TOWN_AUTO_RESOLVE"
	autoSync                    = "GIT_TOWN_AUTO_SYNC"
	bitbucketAppPassword        = "GIT_TOWN_BITBUCKET_APP_PASSWORD"
	bitbucketUserName           = "GIT_TOWN_BITBUCKET_USERNAME"
	branchPrefix                = "GIT_TOWN_BRANCH_PREFIX"
	contributionRegex           = "GIT_TOWN_CONTRIBUTION_REGEX"
	detached                    = "GIT_TOWN_DETACHED"
	devRemote                   = "GIT_TOWN_DEV_REMOTE"
	displayTypes                = "GIT_TOWN_DISPLAY_TYPES"
	dryRun                      = "GIT_TOWN_DRY_RUN"
	featureRegex                = "GIT_TOWN_FEATURE_REGEX"
	forgeType                   = "GIT_TOWN_FORGE_TYPE"
	forgejoToken                = "GIT_TOWN_FORGEJO_TOKEN"
	gitAuthorEmail              = "GIT_AUTHOR_EMAIL"
	gitAuthorName               = "GIT_AUTHOR_NAME"
	gitCommitterEmail           = "GIT_COMMITTER_EMAIL"
	gitCommitterName            = "GIT_COMMITTER_NAME"
	giteaToken                  = "GIT_TOWN_GITEA_TOKEN"
	githubConnectorType         = "GIT_TOWN_GITHUB_CONNECTOR"
	githubToken                 = "GIT_TOWN_GITHUB_TOKEN"
	gitlabConnectorType         = "GIT_TOWN_GITLAB_CONNECTOR"
	gitlabToken                 = "GIT_TOWN_GITLAB_TOKEN"
	ignoreUncommitted           = "GIT_TOWN_IGNORE_UNCOMMITTED"
	interactive                 = "GIT_TOWN_INTERACTIVE"
	mainBranch                  = "GIT_TOWN_MAIN_BRANCH"
	newBranchType               = "GIT_TOWN_NEW_BRANCH_TYPE"
	observedRegex               = "GIT_TOWN_OBSERVED_REGEX"
	offline                     = "GIT_TOWN_OFFLINE"
	order                       = "GIT_TOWN_ORDER"
	originHostname              = "GIT_TOWN_ORIGIN_HOSTNAME"
	perennialBranches           = "GIT_TOWN_PERENNIAL_BRANCHES"
	perennialRegex              = "GIT_TOWN_PERENNIAL_REGEX"
	proposalBreadcrumb          = "GIT_TOWN_PROPOSAL_BREADCRUMB"
	proposalBreadcrumbDirection = "GIT_TOWN_PROPOSAL_BREADCRUMB_DIRECTION"
	pushBranches                = "GIT_TOWN_PUSH_BRANCHES"
	pushHook                    = "GIT_TOWN_PUSH_HOOK"
	shareNewBranches            = "GIT_TOWN_SHARE_NEW_BRANCHES"
	shipDeleteTrackingBranch    = "GIT_TOWN_SHIP_DELETE_TRACKING_BRANCH"
	shipEnterMessage            = "GIT_TOWN_SHIP_ENTER_MESSAGE"
	shipStrategy                = "GIT_TOWN_SHIP_STRATEGY"
	stash                       = "GIT_TOWN_STASH"
	syncFeatureStrategy         = "GIT_TOWN_SYNC_FEATURE_STRATEGY"
	syncPerennialStrategy       = "GIT_TOWN_SYNC_PERENNIAL_STRATEGY"
	syncPrototypeStrategy       = "GIT_TOWN_SYNC_PROTOTYPE_STRATEGY"
	syncTags                    = "GIT_TOWN_SYNC_TAGS"
	syncUpstream                = "GIT_TOWN_SYNC_UPSTREAM"
	term                        = "TERM"
	unknownBranchType           = "GIT_TOWN_UNKNOWN_BRANCH_TYPE"
	verbose                     = "GIT_TOWN_VERBOSE"
	// keep-sorted end
)

func Load(env EnvVars) (configdomain.PartialConfig, error) {
	// keep-sorted start
	autoResolve, errAutoResolve := load(env, autoResolve, gohacks.StrOpt2BoolOpt[configdomain.AutoResolve])
	autoSync, errAutoSync := load(env, autoSync, gohacks.StrOpt2BoolOpt[configdomain.AutoSync])
	branchPrefix, errBranchPrefix := load(env, branchPrefix, configdomain.ParseBranchPrefixOpt)
	browserExecutable, browserEnabled, errBrowser := browserdomain.ParseBrowserOpt(env.GetOpt(Browser))
	contributionRegex, errContribRegex := load(env, contributionRegex, configdomain.ParseContributionRegexOpt)
	detached, errDetached := load(env, detached, gohacks.StrOpt2BoolOpt[configdomain.Detached])
	displayTypesOpt, errDisplayTypes := load(env, displayTypes, configdomain.ParseDisplayTypesOpt)
	dryRun, errDryRun := load(env, dryRun, gohacks.StrOpt2BoolOpt[configdomain.DryRun])
	featureRegex, errFeatureRegex := load(env, featureRegex, configdomain.ParseFeatureRegexOpt)
	forgeType, errForgeType := load(env, forgeType, forgedomain.ParseForgeTypeOpt)
	gitUserEmail := gitdomain.ParseGitUserEmailOpt(env.GetFirstNonEmpty(gitAuthorEmail, gitCommitterEmail))
	gitUserName := gitdomain.ParseGitUserNameOpt(env.GetFirstNonEmpty(gitAuthorName, gitCommitterName))
	githubConnectorType, errGithubConnectorType := load(env, githubConnectorType, forgedomain.ParseGithubConnectorTypeOpt)
	gitlabConnectorType, errGitlabConnectorType := load(env, gitlabConnectorType, forgedomain.ParseGitlabConnectorTypeOpt)
	ignoreUncommitted, errIgnoreUncommitted := load(env, ignoreUncommitted, gohacks.StrOpt2BoolOpt[configdomain.IgnoreUncommitted])
	interactive1, errInteractive1 := load(env, interactive, gohacks.StrOpt2BoolOpt[bool])
	interactive2 := configdomain.NewInteractiveFromEnv(env.GetOpt(term), interactive1)
	newBranchType, errNewBranchType := load(env, newBranchType, configdomain.ParseBranchTypeOpt)
	observedRegex, errObservedRegex := load(env, observedRegex, configdomain.ParseObservedRegexOpt)
	offline, errOffline := load(env, offline, gohacks.StrOpt2BoolOpt[configdomain.Offline])
	order, errOrder := configdomain.ParseOrderOpt(env.GetOpt(order), order)
	perennialRegex, errPerennialRegex := load(env, perennialRegex, configdomain.ParsePerennialRegexOpt)
	proposalBreadcrumb, errProposalBreadcrumb := load(env, proposalBreadcrumb, configdomain.ParseProposalBreadcrumbOpt)
	proposalBreadcrumbDirection, errProposalBreadcrumbDirection := load(env, proposalBreadcrumbDirection, configdomain.ParseProposalBreadcrumbDirectionOpt)
	pushBranches, errPushBranches := load(env, pushBranches, gohacks.StrOpt2BoolOpt[configdomain.PushBranches])
	pushHook, errPushHook := load(env, pushHook, gohacks.StrOpt2BoolOpt[configdomain.PushHook])
	shareNewBranches, errShareNewBranches := load(env, shareNewBranches, configdomain.ParseShareNewBranchesOpt)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := load(env, shipDeleteTrackingBranch, gohacks.StrOpt2BoolOpt[configdomain.ShipDeleteTrackingBranch])
	shipStrategy, errShipStrategy := load(env, shipStrategy, configdomain.ParseShipStrategyOpt)
	stash, errStash := load(env, stash, gohacks.StrOpt2BoolOpt[configdomain.Stash])
	syncFeatureStrategy, errSyncFeatureStrategy := load(env, syncFeatureStrategy, configdomain.ParseSyncFeatureStrategyOpt)
	syncPerennialStrategy, errSyncPerennialStrategy := load(env, syncPerennialStrategy, configdomain.ParseSyncPerennialStrategyOpt)
	syncPrototypeStrategy, errSyncPrototypeStrategy := load(env, syncPrototypeStrategy, configdomain.ParseSyncPrototypeStrategyOpt)
	syncTags, errSyncTags := load(env, syncTags, gohacks.StrOpt2BoolOpt[configdomain.SyncTags])
	syncUpstream, errSyncUpstream := load(env, syncUpstream, gohacks.StrOpt2BoolOpt[configdomain.SyncUpstream])
	unknownBranchType, errUnknownBranchType := load(env, unknownBranchType, configdomain.ParseBranchTypeOpt)
	verbose, errVerbose := load(env, verbose, gohacks.StrOpt2BoolOpt[configdomain.Verbose])
	// keep-sorted end
	err := cmp.Or(
		errAutoResolve,
		errAutoSync,
		errBranchPrefix,
		errBrowser,
		errContribRegex,
		errDetached,
		errDisplayTypes,
		errDryRun,
		errFeatureRegex,
		errForgeType,
		errGithubConnectorType,
		errGitlabConnectorType,
		errIgnoreUncommitted,
		errInteractive1,
		errNewBranchType,
		errObservedRegex,
		errOffline,
		errOrder,
		errPerennialRegex,
		errProposalBreadcrumb,
		errProposalBreadcrumbDirection,
		errPushBranches,
		errPushHook,
		errShareNewBranches,
		errShipDeleteTrackingBranch,
		errShipEnterMessage,
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
		Aliases:                     configdomain.Aliases{}, // aliases aren't loaded from env vars
		AutoResolve:                 autoResolve,
		AutoSync:                    autoSync,
		BitbucketAppPassword:        forgedomain.ParseBitbucketAppPassword(env.GetOpt(bitbucketAppPassword)),
		BitbucketUsername:           forgedomain.ParseBitbucketUsername(env.GetOpt(bitbucketUserName)),
		BranchPrefix:                branchPrefix,
		BranchTypeOverrides:         configdomain.BranchTypeOverrides{}, // not loaded from env vars
		BrowserEnabled:              browserEnabled,
		BrowserExecutable:           browserExecutable,
		ForgejoToken:                forgedomain.ParseForgejoToken(env.GetOpt(forgejoToken)),
		ContributionRegex:           contributionRegex,
		Detached:                    detached,
		DevRemote:                   gitdomain.NewRemote(env.GetOpt(devRemote)),
		DisplayTypes:                displayTypesOpt,
		DryRun:                      dryRun,
		FeatureRegex:                featureRegex,
		ForgeType:                   forgeType,
		GithubConnectorType:         githubConnectorType,
		GithubToken:                 forgedomain.ParseGithubToken(env.GetOpt(githubToken, "GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")),
		GitlabConnectorType:         gitlabConnectorType,
		GitlabToken:                 forgedomain.ParseGitlabToken(env.GetOpt(gitlabToken)),
		GitUserEmail:                gitUserEmail,
		GitUserName:                 gitUserName,
		GiteaToken:                  forgedomain.ParseGiteaToken(env.GetOpt(giteaToken)),
		HostingOriginHostname:       configdomain.ParseHostingOriginHostname(env.GetOpt(originHostname)),
		IgnoreUncommitted:           ignoreUncommitted,
		Interactive:                 interactive2,
		Lineage:                     configdomain.NewLineage(), // not loaded from env vars
		MainBranch:                  gitdomain.LocalBranchNameOpt(env.Get(mainBranch)),
		NewBranchType:               configdomain.NewBranchTypeOpt(newBranchType),
		ObservedRegex:               observedRegex,
		Offline:                     offline,
		Order:                       order,
		PerennialBranches:           gitdomain.ParseLocalBranchNames(env.GetOpt(perennialBranches)),
		PerennialRegex:              perennialRegex,
		ProposalBreadcrumb:          proposalBreadcrumb,
		ProposalBreadcrumbDirection: proposalBreadcrumbDirection,
		PushBranches:                pushBranches,
		PushHook:                    pushHook,
		ShareNewBranches:            shareNewBranches,
		ShipDeleteTrackingBranch:    shipDeleteTrackingBranch,
		ShipEnterMessage:            shipEnterMessage,
		ShipStrategy:                shipStrategy,
		Stash:                       stash,
		SyncFeatureStrategy:         syncFeatureStrategy,
		SyncPerennialStrategy:       syncPerennialStrategy,
		SyncPrototypeStrategy:       syncPrototypeStrategy,
		SyncTags:                    syncTags,
		SyncUpstream:                syncUpstream,
		UnknownBranchType:           configdomain.UnknownBranchTypeOpt(unknownBranchType),
		Verbose:                     verbose,
	}, err
}

func load[T any](env EnvVars, varName string, parser func(stringss.Trimmed, string) (T, error)) (T, error) { //nolint:ireturn
	return parser(env.Get(varName), varName)
}
