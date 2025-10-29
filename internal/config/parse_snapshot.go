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

// provides the branch type overrides stored in the given Git metadata snapshot
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
		branchTypeOpt, err := configdomain.ParseBranchType(value)
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
	autoResolve, errAutoResolve := load(snapshot, configdomain.KeyAutoResolve, gohacks.ParseBoolOpt[configdomain.AutoResolve], ignoreUnknown, None[configdomain.AutoResolve]())
	autoSync, errAutoSync := load(snapshot, configdomain.KeyAutoSync, gohacks.ParseBoolOpt[configdomain.AutoSync], ignoreUnknown, None[configdomain.AutoSync]())
	branchTypeOverrides, errBranchTypeOverride := NewBranchTypeOverridesInSnapshot(snapshot, ignoreUnknown, runner)
	contributionRegex, errContributionRegex := load(snapshot, configdomain.KeyContributionRegex, configdomain.ParseContributionRegex, ignoreUnknown, None[configdomain.ContributionRegex]())
	detached, errDetached := load(snapshot, configdomain.KeyDetached, gohacks.ParseBoolOpt[configdomain.Detached], ignoreUnknown, None[configdomain.Detached]())
	displayTypes, errDisplayTypes := load(snapshot, configdomain.KeyDisplayTypes, configdomain.ParseDisplayTypes, ignoreUnknown, None[configdomain.DisplayTypes]())
	featureRegex, errFeatureRegex := load(snapshot, configdomain.KeyFeatureRegex, configdomain.ParseFeatureRegex, ignoreUnknown, None[configdomain.FeatureRegex]())
	forgeType, errForgeType := load(snapshot, configdomain.KeyForgeType, forgedomain.ParseForgeType, ignoreUnknown, None[forgedomain.ForgeType]())
	githubConnectorType, errGitHubConnectorType := load(snapshot, configdomain.KeyGitHubConnectorType, forgedomain.ParseGitHubConnectorType, ignoreUnknown, None[forgedomain.GitHubConnectorType]())
	gitlabConnectorType, errGitLabConnectorType := load(snapshot, configdomain.KeyGitLabConnectorType, forgedomain.ParseGitLabConnectorType, ignoreUnknown, None[forgedomain.GitLabConnectorType]())
	lineage, errLineage := NewLineageFromSnapshot(snapshot, updateOutdated, runner)
	newBranchTypeValue, errNewBranchType := load(snapshot, configdomain.KeyNewBranchType, configdomain.ParseBranchType, ignoreUnknown, None[configdomain.BranchType]())
	newBranchType := configdomain.NewBranchTypeOpt(newBranchTypeValue)
	observedRegex, errObservedRegex := load(snapshot, configdomain.KeyObservedRegex, configdomain.ParseObservedRegex, ignoreUnknown, None[configdomain.ObservedRegex]())
	order, errOrder := load(snapshot, configdomain.KeyOrder, configdomain.ParseOrder, ignoreUnknown, None[configdomain.Order]())
	offline, errOffline := load(snapshot, configdomain.KeyOffline, gohacks.ParseBoolOpt[configdomain.Offline], ignoreUnknown, None[configdomain.Offline]())
	perennialRegex, errPerennialRegex := load(snapshot, configdomain.KeyPerennialRegex, configdomain.ParsePerennialRegex, ignoreUnknown, None[configdomain.PerennialRegex]())
	proposalsShowLineage, errProposalsShowLineage := load(snapshot, configdomain.KeyProposalsShowLineage, forgedomain.ParseProposalsShowLineage, ignoreUnknown, None[forgedomain.ProposalsShowLineage]())
	pushBranches, errPushBranches := load(snapshot, configdomain.KeyPushBranches, gohacks.ParseBoolOpt[configdomain.PushBranches], ignoreUnknown, None[configdomain.PushBranches]())
	pushHook, errPushHook := load(snapshot, configdomain.KeyPushHook, gohacks.ParseBoolOpt[configdomain.PushHook], ignoreUnknown, None[configdomain.PushHook]())
	shareNewBranches, errShareNewBranches := load(snapshot, configdomain.KeyShareNewBranches, configdomain.ParseShareNewBranches)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := load(snapshot, configdomain.KeyShipDeleteTrackingBranch, gohacks.ParseBoolOpt[configdomain.ShipDeleteTrackingBranch])
	shipStrategy, errShipStrategy := load(snapshot, configdomain.KeyShipStrategy, configdomain.ParseShipStrategy)
	stash, errStash := load(snapshot, configdomain.KeyStash, gohacks.ParseBoolOpt[configdomain.Stash])
	syncFeatureStrategy, errSyncFeatureStrategy := load(snapshot, configdomain.KeySyncFeatureStrategy, configdomain.ParseSyncFeatureStrategy)
	syncPerennialStrategy, errSyncPerennialStrategy := load(snapshot, configdomain.KeySyncPerennialStrategy, configdomain.ParseSyncPerennialStrategy)
	syncPrototypeStrategy, errSyncPrototypeStrategy := load(snapshot, configdomain.KeySyncPrototypeStrategy, configdomain.ParseSyncPrototypeStrategy)
	syncTags, errSyncTags := load(snapshot, configdomain.KeySyncTags, gohacks.ParseBoolOpt[configdomain.SyncTags])
	syncUpstream, errSyncUpstream := load(snapshot, configdomain.KeySyncUpstream, gohacks.ParseBoolOpt[configdomain.SyncUpstream])
	unknownBranchTypeValue, errUnknownBranchType := load(snapshot, configdomain.KeyUnknownBranchType, configdomain.ParseBranchType)
	unknownBranchType := configdomain.UnknownBranchTypeOpt(unknownBranchTypeValue)
	err := cmp.Or(
		errAutoResolve,
		errAutoSync,
		errBranchTypeOverride,
		errContributionRegex,
		errDetached,
		errDisplayTypes,
		errFeatureRegex,
		errForgeType,
		errGitHubConnectorType,
		errGitLabConnectorType,
		errLineage,
		errNewBranchType,
		errObservedRegex,
		errOrder,
		errOffline,
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
	)
	return configdomain.PartialConfig{
		Aliases:                  snapshot.Aliases(),
		AutoResolve:              autoResolve,
		AutoSync:                 autoSync,
		BitbucketAppPassword:     forgedomain.ParseBitbucketAppPassword(snapshot[configdomain.KeyBitbucketAppPassword]),
		BitbucketUsername:        forgedomain.ParseBitbucketUsername(snapshot[configdomain.KeyBitbucketUsername]),
		BranchTypeOverrides:      branchTypeOverrides,
		ForgejoToken:             forgedomain.ParseForgejoToken(snapshot[configdomain.KeyForgejoToken]),
		ContributionRegex:        contributionRegex,
		Detached:                 detached,
		DevRemote:                gitdomain.NewRemote(snapshot[configdomain.KeyDevRemote]),
		DisplayTypes:             displayTypes,
		DryRun:                   None[configdomain.DryRun](),
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubConnectorType:      githubConnectorType,
		GitHubToken:              forgedomain.ParseGitHubToken(snapshot[configdomain.KeyGitHubToken]),
		GitLabConnectorType:      gitlabConnectorType,
		GitLabToken:              forgedomain.ParseGitLabToken(snapshot[configdomain.KeyGitLabToken]),
		GitUserEmail:             gitdomain.ParseGitUserEmail(snapshot[configdomain.KeyGitUserEmail]),
		GitUserName:              gitdomain.ParseGitUserName(snapshot[configdomain.KeyGitUserName]),
		GiteaToken:               forgedomain.ParseGiteaToken(snapshot[configdomain.KeyGiteaToken]),
		HostingOriginHostname:    configdomain.ParseHostingOriginHostname(snapshot[configdomain.KeyHostingOriginHostname]),
		Lineage:                  lineage,
		MainBranch:               gitdomain.NewLocalBranchNameOption(snapshot[configdomain.KeyMainBranch]),
		NewBranchType:            newBranchType,
		ObservedRegex:            observedRegex,
		Order:                    order,
		Offline:                  offline,
		PerennialBranches:        gitdomain.ParseLocalBranchNames(snapshot[configdomain.KeyPerennialBranches]),
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
		UnknownBranchType:        unknownBranchType,
		Verbose:                  None[configdomain.Verbose](),
	}, err
}

func load[T any](snapshot configdomain.SingleSnapshot, key configdomain.Key, parseFunc func(string, string) (T, error), ignoreUnknown bool, defaults T) (T, error) {
	value, err := parseFunc(snapshot[key], key.String())
	if err != nil {
		if ignoreUnknown {
			return defaults, nil
		}
		return defaults, err
	}
	return value, nil
}
