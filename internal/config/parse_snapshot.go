package config

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// NewBranchTypeOverridesInSnapshot provides the branch type overrides stored in the given Git metadata snapshot.
func NewBranchTypeOverridesInSnapshot(snapshot configdomain.SingleSnapshot, ignoreUnknown bool, runner subshelldomain.Runner) (configdomain.BranchTypeOverrides, error) {
	result := configdomain.BranchTypeOverrides{}
	for key, value := range snapshot { // okay to iterate the map in random order because we assign to a new map
		key, isBranchTypeKey := configdomain.ParseBranchTypeOverrideKey(key).Get()
		if !isBranchTypeKey {
			continue
		}
		branch := key.Branch()
		if branch == "" {
			// empty branch name --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigBranchTypeOverrideEmpty))
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			// empty branch type values are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigBranchTypeOverrideEmpty))
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		branchTypeOpt, err := configdomain.ParseBranchType(value, key.String())
		if err != nil {
			if ignoreUnknown {
				fmt.Printf("Ignoring unknown branch type override for %q: %s\n", branch, value)
			} else {
				return result, err
			}
		}
		if branchType, hasBranchType := branchTypeOpt.Get(); hasBranchType {
			result[branch] = branchType
		}
	}
	return result, nil
}

func NewLineageFromSnapshot(snapshot configdomain.SingleSnapshot, updateOutdated bool, runner subshelldomain.Runner) (configdomain.Lineage, error) {
	result := configdomain.NewLineage()
	for key, value := range snapshot.LineageEntries() { // okay to iterate the map in random order because we assign to a new map
		child := key.ChildBranch()
		if child == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigLineageEmptyChild))
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigLineageEmptyChild))
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		if updateOutdated && child.String() == value {
			fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.ConfigLineageParentIsChild, child)))
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
		}
		parent := gitdomain.NewLocalBranchName(value)
		result = result.Set(child, parent)
	}
	return result, nil
}

func NewPartialConfigFromSnapshot(snapshot configdomain.SingleSnapshot, updateOutdated bool, ignoreUnknown bool, runner subshelldomain.Runner) (configdomain.PartialConfig, error) {
	// TODO: add keep-sorted to all blocks in this function
	autoResolve, errAutoResolve := load(snapshot, configdomain.KeyAutoResolve, gohacks.ParseBoolOpt[configdomain.AutoResolve], ignoreUnknown)
	autoSync, errAutoSync := load(snapshot, configdomain.KeyAutoSync, gohacks.ParseBoolOpt[configdomain.AutoSync], ignoreUnknown)
	branchPrefix, errBranchPrefix := load(snapshot, configdomain.KeyBranchPrefix, configdomain.ParseBranchPrefix, ignoreUnknown)
	branchTypeOverrides, errBranchTypeOverride := NewBranchTypeOverridesInSnapshot(snapshot, ignoreUnknown, runner)
	browser, errBrowser := load(snapshot, configdomain.KeyBrowser, configdomain.ParseBrowser, ignoreUnknown)
	contributionRegex, errContributionRegex := load(snapshot, configdomain.KeyContributionRegex, configdomain.ParseContributionRegex, ignoreUnknown)
	detached, errDetached := load(snapshot, configdomain.KeyDetached, gohacks.ParseBoolOpt[configdomain.Detached], ignoreUnknown)
	displayTypes, errDisplayTypes := load(snapshot, configdomain.KeyDisplayTypes, configdomain.ParseDisplayTypes, ignoreUnknown)
	featureRegex, errFeatureRegex := load(snapshot, configdomain.KeyFeatureRegex, configdomain.ParseFeatureRegex, ignoreUnknown)
	forgeType, errForgeType := load(snapshot, configdomain.KeyForgeType, forgedomain.ParseForgeType, ignoreUnknown)
	githubConnectorType, errGithubConnectorType := load(snapshot, configdomain.KeyGithubConnectorType, forgedomain.ParseGithubConnectorType, ignoreUnknown)
	gitlabConnectorType, errGitlabConnectorType := load(snapshot, configdomain.KeyGitlabConnectorType, forgedomain.ParseGitlabConnectorType, ignoreUnknown)
	ignoreUncommitted, errIgnoreUncommitted := load(snapshot, configdomain.KeyIgnoreUncommitted, gohacks.ParseBoolOpt[configdomain.IgnoreUncommitted], ignoreUnknown)
	lineage, errLineage := NewLineageFromSnapshot(snapshot, updateOutdated, runner)
	newBranchTypeValue, errNewBranchType := load(snapshot, configdomain.KeyNewBranchType, configdomain.ParseBranchType, ignoreUnknown)
	newBranchType := configdomain.NewBranchTypeOpt(newBranchTypeValue)
	observedRegex, errObservedRegex := load(snapshot, configdomain.KeyObservedRegex, configdomain.ParseObservedRegex, ignoreUnknown)
	order, errOrder := load(snapshot, configdomain.KeyOrder, configdomain.ParseOrder, ignoreUnknown)
	offline, errOffline := load(snapshot, configdomain.KeyOffline, gohacks.ParseBoolOpt[configdomain.Offline], ignoreUnknown)
	perennialRegex, errPerennialRegex := load(snapshot, configdomain.KeyPerennialRegex, configdomain.ParsePerennialRegex, ignoreUnknown)
	proposalBreadcrumb, errProposalBreadcrumb := load(snapshot, configdomain.KeyProposalBreadcrumb, configdomain.ParseProposalBreadcrumb, ignoreUnknown)
	proposalBreadcrumbDirection, errProposalBreadcrumbDirection := load(snapshot, configdomain.KeyProposalBreadcrumbDirection, configdomain.ParseProposalBreadcrumbDirection, ignoreUnknown)
	pushBranches, errPushBranches := load(snapshot, configdomain.KeyPushBranches, gohacks.ParseBoolOpt[configdomain.PushBranches], ignoreUnknown)
	pushHook, errPushHook := load(snapshot, configdomain.KeyPushHook, gohacks.ParseBoolOpt[configdomain.PushHook], ignoreUnknown)
	shareNewBranches, errShareNewBranches := load(snapshot, configdomain.KeyShareNewBranches, configdomain.ParseShareNewBranches, ignoreUnknown)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := load(snapshot, configdomain.KeyShipDeleteTrackingBranch, gohacks.ParseBoolOpt[configdomain.ShipDeleteTrackingBranch], ignoreUnknown)
	shipStrategy, errShipStrategy := load(snapshot, configdomain.KeyShipStrategy, configdomain.ParseShipStrategy, ignoreUnknown)
	stash, errStash := load(snapshot, configdomain.KeyStash, gohacks.ParseBoolOpt[configdomain.Stash], ignoreUnknown)
	syncFeatureStrategy, errSyncFeatureStrategy := load(snapshot, configdomain.KeySyncFeatureStrategy, configdomain.ParseSyncFeatureStrategy, ignoreUnknown)
	syncPerennialStrategy, errSyncPerennialStrategy := load(snapshot, configdomain.KeySyncPerennialStrategy, configdomain.ParseSyncPerennialStrategy, ignoreUnknown)
	syncPrototypeStrategy, errSyncPrototypeStrategy := load(snapshot, configdomain.KeySyncPrototypeStrategy, configdomain.ParseSyncPrototypeStrategy, ignoreUnknown)
	syncTags, errSyncTags := load(snapshot, configdomain.KeySyncTags, gohacks.ParseBoolOpt[configdomain.SyncTags], ignoreUnknown)
	syncUpstream, errSyncUpstream := load(snapshot, configdomain.KeySyncUpstream, gohacks.ParseBoolOpt[configdomain.SyncUpstream], ignoreUnknown)
	unknownBranchTypeValue, errUnknownBranchType := load(snapshot, configdomain.KeyUnknownBranchType, configdomain.ParseBranchType, ignoreUnknown)
	unknownBranchType := configdomain.UnknownBranchTypeOpt(unknownBranchTypeValue)
	err := cmp.Or(
		errAutoResolve,
		errAutoSync,
		errBranchPrefix,
		errBranchTypeOverride,
		errBrowser,
		errContributionRegex,
		errDetached,
		errDisplayTypes,
		errFeatureRegex,
		errForgeType,
		errGithubConnectorType,
		errGitlabConnectorType,
		errIgnoreUncommitted,
		errLineage,
		errNewBranchType,
		errObservedRegex,
		errOrder,
		errOffline,
		errPerennialRegex,
		errProposalBreadcrumb,
		errProposalBreadcrumbDirection,
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
	)
	return configdomain.PartialConfig{
		Aliases:                     snapshot.Aliases(),
		AutoResolve:                 autoResolve,
		AutoSync:                    autoSync,
		BitbucketAppPassword:        forgedomain.ParseBitbucketAppPassword(snapshot[configdomain.KeyBitbucketAppPassword]),
		BitbucketUsername:           forgedomain.ParseBitbucketUsername(snapshot[configdomain.KeyBitbucketUsername]),
		BranchPrefix:                branchPrefix,
		BranchTypeOverrides:         branchTypeOverrides,
		Browser:                     browser,
		ForgejoToken:                forgedomain.ParseForgejoToken(snapshot[configdomain.KeyForgejoToken]),
		ContributionRegex:           contributionRegex,
		Detached:                    detached,
		DevRemote:                   gitdomain.NewRemote(snapshot[configdomain.KeyDevRemote]),
		DisplayTypes:                displayTypes,
		DryRun:                      None[configdomain.DryRun](),
		FeatureRegex:                featureRegex,
		ForgeType:                   forgeType,
		GithubConnectorType:         githubConnectorType,
		GithubToken:                 forgedomain.ParseGithubToken(snapshot[configdomain.KeyGithubToken]),
		GitlabConnectorType:         gitlabConnectorType,
		GitlabToken:                 forgedomain.ParseGitlabToken(snapshot[configdomain.KeyGitlabToken]),
		GitUserEmail:                gitdomain.ParseGitUserEmail(snapshot[configdomain.KeyGitUserEmail]),
		GitUserName:                 gitdomain.ParseGitUserName(snapshot[configdomain.KeyGitUserName]),
		GiteaToken:                  forgedomain.ParseGiteaToken(snapshot[configdomain.KeyGiteaToken]),
		HostingOriginHostname:       configdomain.ParseHostingOriginHostname(snapshot[configdomain.KeyHostingOriginHostname]),
		IgnoreUncommitted:           ignoreUncommitted,
		Lineage:                     lineage,
		MainBranch:                  gitdomain.NewLocalBranchNameOption(snapshot[configdomain.KeyMainBranch]),
		NewBranchType:               newBranchType,
		ObservedRegex:               observedRegex,
		Order:                       order,
		Offline:                     offline,
		PerennialBranches:           gitdomain.ParseLocalBranchNames(snapshot[configdomain.KeyPerennialBranches]),
		PerennialRegex:              perennialRegex,
		ProposalBreadcrumb:          proposalBreadcrumb,
		ProposalBreadcrumbDirection: proposalBreadcrumbDirection,
		PushBranches:                pushBranches,
		PushHook:                    pushHook,
		ShareNewBranches:            shareNewBranches,
		ShipDeleteTrackingBranch:    shipDeleteTrackingBranch,
		ShipStrategy:                shipStrategy,
		Stash:                       stash,
		SyncFeatureStrategy:         syncFeatureStrategy,
		SyncPerennialStrategy:       syncPerennialStrategy,
		SyncPrototypeStrategy:       syncPrototypeStrategy,
		SyncTags:                    syncTags,
		SyncUpstream:                syncUpstream,
		UnknownBranchType:           unknownBranchType,
		Verbose:                     None[configdomain.Verbose](),
	}, err
}

func load[T any](snapshot configdomain.SingleSnapshot, key configdomain.Key, parseFunc func(string, string) (T, error), ignoreUnknown bool) (T, error) { //nolint:ireturn
	valueStr, has := snapshot[key]
	if !has {
		var zero T
		return zero, nil
	}
	value, err := parseFunc(valueStr, key.String())
	if err != nil {
		if ignoreUnknown {
			fmt.Printf("Ignoring invalid value for %q: %q\n", key, valueStr)
			var zero T
			return zero, nil
		}
		var zero T
		return zero, err
	}
	return value, nil
}
