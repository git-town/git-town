package config

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/browser/browserdomain"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/config/gitconfig"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks"
	"github.com/git-town/git-town/v23/internal/messages"
	"github.com/git-town/git-town/v23/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v23/pkg/colors"
	. "github.com/git-town/git-town/v23/pkg/prelude"
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
	browserStr, hasBrowser := snapshot[configdomain.KeyBrowser]
	// keep-sorted start
	autoResolve, errAutoResolve := loadErr(snapshot, configdomain.KeyAutoResolve, gohacks.Str2BoolOpt[configdomain.AutoResolve], ignoreUnknown)
	autoSync, errAutoSync := loadErr(snapshot, configdomain.KeyAutoSync, gohacks.Str2BoolOpt[configdomain.AutoSync], ignoreUnknown)
	branchPrefix, errBranchPrefix := loadErr(snapshot, configdomain.KeyBranchPrefix, configdomain.ParseBranchPrefix, ignoreUnknown)
	branchTypeOverrides, errBranchTypeOverride := NewBranchTypeOverridesInSnapshot(snapshot, ignoreUnknown, runner)
	browserExecutable, browserEnabled, errBrowser := browserdomain.ParseBrowserOpt(NewOptionIfExists(browserStr, hasBrowser))
	contributionRegex, errContributionRegex := loadErr(snapshot, configdomain.KeyContributionRegex, configdomain.ParseContributionRegex, ignoreUnknown)
	detached, errDetached := loadErr(snapshot, configdomain.KeyDetached, gohacks.Str2BoolOpt[configdomain.Detached], ignoreUnknown)
	displayTypes, errDisplayTypes := loadErr(snapshot, configdomain.KeyDisplayTypes, configdomain.ParseDisplayTypes, ignoreUnknown)
	featureRegex, errFeatureRegex := loadErr(snapshot, configdomain.KeyFeatureRegex, configdomain.ParseFeatureRegex, ignoreUnknown)
	forgeType, errForgeType := loadErr(snapshot, configdomain.KeyForgeType, forgedomain.ParseForgeType, ignoreUnknown)
	githubConnectorType, errGithubConnectorType := loadErr(snapshot, configdomain.KeyGithubConnectorType, forgedomain.ParseGithubConnectorType, ignoreUnknown)
	gitlabConnectorType, errGitlabConnectorType := loadErr(snapshot, configdomain.KeyGitlabConnectorType, forgedomain.ParseGitlabConnectorType, ignoreUnknown)
	ignoreUncommitted, errIgnoreUncommitted := loadErr(snapshot, configdomain.KeyIgnoreUncommitted, gohacks.Str2BoolOpt[configdomain.IgnoreUncommitted], ignoreUnknown)
	interactive, errInteractive := loadErr(snapshot, configdomain.KeyInteractive, configdomain.NewInteractiveFromSnapshot, ignoreUnknown)
	lineage, errLineage := NewLineageFromSnapshot(snapshot, updateOutdated, runner)
	newBranchType1, errNewBranchType := loadErr(snapshot, configdomain.KeyNewBranchType, configdomain.ParseBranchType, ignoreUnknown)
	newBranchType2 := configdomain.NewBranchTypeOpt(newBranchType1)
	observedRegex, errObservedRegex := loadErr(snapshot, configdomain.KeyObservedRegex, configdomain.ParseObservedRegex, ignoreUnknown)
	offline, errOffline := loadErr(snapshot, configdomain.KeyOffline, gohacks.Str2BoolOpt[configdomain.Offline], ignoreUnknown)
	order, errOrder := loadErr(snapshot, configdomain.KeyOrder, configdomain.ParseOrder, ignoreUnknown)
	perennialRegex, errPerennialRegex := loadErr(snapshot, configdomain.KeyPerennialRegex, configdomain.ParsePerennialRegex, ignoreUnknown)
	proposalBreadcrumb, errProposalBreadcrumb := loadErr(snapshot, configdomain.KeyProposalBreadcrumb, configdomain.ParseProposalBreadcrumb, ignoreUnknown)
	proposalBreadcrumbDirection, errProposalBreadcrumbDirection := loadErr(snapshot, configdomain.KeyProposalBreadcrumbDirection, configdomain.ParseProposalBreadcrumbDirection, ignoreUnknown)
	pushBranches, errPushBranches := loadErr(snapshot, configdomain.KeyPushBranches, gohacks.Str2BoolOpt[configdomain.PushBranches], ignoreUnknown)
	pushHook, errPushHook := loadErr(snapshot, configdomain.KeyPushHook, gohacks.Str2BoolOpt[configdomain.PushHook], ignoreUnknown)
	shareNewBranches, errShareNewBranches := loadErr(snapshot, configdomain.KeyShareNewBranches, configdomain.ParseShareNewBranches, ignoreUnknown)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := loadErr(snapshot, configdomain.KeyShipDeleteTrackingBranch, gohacks.Str2BoolOpt[configdomain.ShipDeleteTrackingBranch], ignoreUnknown)
	shipStrategy, errShipStrategy := loadErr(snapshot, configdomain.KeyShipStrategy, configdomain.ParseShipStrategy, ignoreUnknown)
	stash, errStash := loadErr(snapshot, configdomain.KeyStash, gohacks.Str2BoolOpt[configdomain.Stash], ignoreUnknown)
	syncFeatureStrategy, errSyncFeatureStrategy := loadErr(snapshot, configdomain.KeySyncFeatureStrategy, configdomain.ParseSyncFeatureStrategy, ignoreUnknown)
	syncPerennialStrategy, errSyncPerennialStrategy := loadErr(snapshot, configdomain.KeySyncPerennialStrategy, configdomain.ParseSyncPerennialStrategy, ignoreUnknown)
	syncPrototypeStrategy, errSyncPrototypeStrategy := loadErr(snapshot, configdomain.KeySyncPrototypeStrategy, configdomain.ParseSyncPrototypeStrategy, ignoreUnknown)
	syncTags, errSyncTags := loadErr(snapshot, configdomain.KeySyncTags, gohacks.Str2BoolOpt[configdomain.SyncTags], ignoreUnknown)
	syncUpstream, errSyncUpstream := loadErr(snapshot, configdomain.KeySyncUpstream, gohacks.Str2BoolOpt[configdomain.SyncUpstream], ignoreUnknown)
	unknownBranchType1, errUnknownBranchType := loadErr(snapshot, configdomain.KeyUnknownBranchType, configdomain.ParseBranchType, ignoreUnknown)
	unknownBranchType2 := configdomain.UnknownBranchTypeOpt(unknownBranchType1)
	// keep-sorted end
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
		errInteractive,
		errLineage,
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
		BitbucketAppPassword:        load(snapshot, configdomain.KeyBitbucketAppPassword, forgedomain.ParseBitbucketAppPassword),
		BitbucketUsername:           load(snapshot, configdomain.KeyBitbucketUsername, forgedomain.ParseBitbucketUsername),
		BranchPrefix:                branchPrefix,
		BranchTypeOverrides:         branchTypeOverrides,
		BrowserExecutable:           browserExecutable,
		BrowserEnabled:              browserEnabled,
		ForgejoToken:                load(snapshot, configdomain.KeyForgejoToken, forgedomain.ParseForgejoToken),
		ContributionRegex:           contributionRegex,
		Detached:                    detached,
		DevRemote:                   load(snapshot, configdomain.KeyDevRemote, gitdomain.NewRemote),
		DisplayTypes:                displayTypes,
		DryRun:                      None[configdomain.DryRun](),
		FeatureRegex:                featureRegex,
		ForgeType:                   forgeType,
		GithubConnectorType:         githubConnectorType,
		GithubToken:                 load(snapshot, configdomain.KeyGithubToken, forgedomain.ParseGithubToken),
		GitlabConnectorType:         gitlabConnectorType,
		GitlabToken:                 load(snapshot, configdomain.KeyGitlabToken, forgedomain.ParseGitlabToken),
		GitUserEmail:                gitdomain.ParseGitUserEmail(snapshot[configdomain.KeyGitUserEmail]),
		GitUserName:                 gitdomain.ParseGitUserName(snapshot[configdomain.KeyGitUserName]),
		GiteaToken:                  load(snapshot, configdomain.KeyGiteaToken, forgedomain.ParseGiteaToken),
		HostingOriginHostname:       load(snapshot, configdomain.KeyHostingOriginHostname, configdomain.ParseHostingOriginHostname),
		IgnoreUncommitted:           ignoreUncommitted,
		Interactive:                 interactive,
		Lineage:                     lineage,
		MainBranch:                  load(snapshot, configdomain.KeyMainBranch, gitdomain.NewLocalBranchNameOption),
		NewBranchType:               newBranchType2,
		ObservedRegex:               observedRegex,
		Order:                       order,
		Offline:                     offline,
		PerennialBranches:           load(snapshot, configdomain.KeyPerennialBranches, gitdomain.ParseLocalBranchNames),
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
		UnknownBranchType:           unknownBranchType2,
		Verbose:                     None[configdomain.Verbose](),
	}, err
}

func load[T any](snapshot configdomain.SingleSnapshot, key configdomain.Key, parseFunc func(Option[string]) T) T { //nolint:ireturn
	valueStr := snapshot.GetOpt(key)
	return parseFunc(valueStr)
}

func loadErr[T any](snapshot configdomain.SingleSnapshot, key configdomain.Key, parseFunc func(string, string) (T, error), ignoreUnknown bool) (T, error) { //nolint:ireturn
	valueStr, has := snapshot[key]
	if !has {
		var zero T
		return zero, nil
	}
	value, err := parseFunc(valueStr, key.String())
	if err != nil {
		var zero T
		if ignoreUnknown {
			fmt.Printf("Ignoring invalid value for %q: %q\n", key, valueStr)
			return zero, nil
		}
		return zero, err
	}
	return value, nil
}
