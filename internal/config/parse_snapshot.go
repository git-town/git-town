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
func NewBranchTypeOverridesInSnapshot(snapshot configdomain.SingleSnapshot, runner subshelldomain.Runner) (configdomain.BranchTypeOverrides, error) {
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
			return result, err
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

func loadField[T any](snapshot configdomain.SingleSnapshot, key configdomain.Key, parseFunc func(string, string) (T, error)) (T, error) {
	return parseFunc(snapshot[key], key.String())
}

func NewPartialConfigFromSnapshot(snapshot configdomain.SingleSnapshot, updateOutdated bool, runner subshelldomain.Runner) (configdomain.PartialConfig, error) {
	autoResolve, errAutoResolve := loadField(snapshot, configdomain.KeyAutoResolve, gohacks.ParseBoolOpt[configdomain.AutoResolve])
	autoSync, errAutoSync := loadField(snapshot, configdomain.KeyAutoSync, gohacks.ParseBoolOpt[configdomain.AutoSync])
	branchTypeOverrides, errBranchTypeOverride := NewBranchTypeOverridesInSnapshot(snapshot, runner)
	contributionRegex, errContributionRegex := loadField(snapshot, configdomain.KeyContributionRegex, configdomain.ParseContributionRegex)
	detached, errDetached := loadField(snapshot, configdomain.KeyDetached, gohacks.ParseBoolOpt[configdomain.Detached])
	displayTypes, errDisplayTypes := loadField(snapshot, configdomain.KeyDisplayTypes, configdomain.ParseDisplayTypes)
	featureRegex, errFeatureRegex := loadField(snapshot, configdomain.KeyFeatureRegex, configdomain.ParseFeatureRegex)
	forgeType, errForgeType := loadField(snapshot, configdomain.KeyForgeType, forgedomain.ParseForgeType)
	githubConnectorType, errGitHubConnectorType := loadField(snapshot, configdomain.KeyGitHubConnectorType, forgedomain.ParseGitHubConnectorType)
	gitlabConnectorType, errGitLabConnectorType := loadField(snapshot, configdomain.KeyGitLabConnectorType, forgedomain.ParseGitLabConnectorType)
	lineage, errLineage := NewLineageFromSnapshot(snapshot, updateOutdated, runner)
	newBranchTypeValue, errNewBranchType := loadField(snapshot, configdomain.KeyNewBranchType, configdomain.ParseBranchType)
	newBranchType := configdomain.NewBranchTypeOpt(newBranchTypeValue)
	observedRegex, errObservedRegex := loadField(snapshot, configdomain.KeyObservedRegex, configdomain.ParseObservedRegex)
	order, errOrder := configdomain.ParseOrder(snapshot[configdomain.KeyOrder], configdomain.KeyOrder)
	offline, errOffline := gohacks.ParseBoolOpt[configdomain.Offline](snapshot[configdomain.KeyOffline], configdomain.KeyOffline.String())
	perennialRegex, errPerennialRegex := configdomain.ParsePerennialRegex(snapshot[configdomain.KeyPerennialRegex])
	proposalsShowLineage, errProposalsShowLineage := forgedomain.ParseProposalsShowLineage(snapshot[configdomain.KeyProposalsShowLineage])
	pushBranches, errPushBranches := gohacks.ParseBoolOpt[configdomain.PushBranches](snapshot[configdomain.KeyPushBranches], configdomain.KeyPushBranches.String())
	pushHook, errPushHook := gohacks.ParseBoolOpt[configdomain.PushHook](snapshot[configdomain.KeyPushHook], configdomain.KeyPushHook.String())
	shareNewBranches, errShareNewBranches := configdomain.ParseShareNewBranches(snapshot[configdomain.KeyShareNewBranches], configdomain.KeyShareNewBranches)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := gohacks.ParseBoolOpt[configdomain.ShipDeleteTrackingBranch](snapshot[configdomain.KeyShipDeleteTrackingBranch], configdomain.KeyShipDeleteTrackingBranch.String())
	shipStrategy, errShipStrategy := configdomain.ParseShipStrategy(snapshot[configdomain.KeyShipStrategy])
	stash, errStash := gohacks.ParseBoolOpt[configdomain.Stash](snapshot[configdomain.KeyStash], configdomain.KeyStash.String())
	syncFeatureStrategy, errSyncFeatureStrategy := configdomain.ParseSyncFeatureStrategy(snapshot[configdomain.KeySyncFeatureStrategy])
	syncPerennialStrategy, errSyncPerennialStrategy := configdomain.ParseSyncPerennialStrategy(snapshot[configdomain.KeySyncPerennialStrategy])
	syncPrototypeStrategy, errSyncPrototypeStrategy := configdomain.ParseSyncPrototypeStrategy(snapshot[configdomain.KeySyncPrototypeStrategy])
	syncTags, errSyncTags := gohacks.ParseBoolOpt[configdomain.SyncTags](snapshot[configdomain.KeySyncTags], configdomain.KeySyncTags.String())
	syncUpstream, errSyncUpstream := gohacks.ParseBoolOpt[configdomain.SyncUpstream](snapshot[configdomain.KeySyncUpstream], configdomain.KeySyncUpstream.String())
	unknownBranchTypeValue, errUnknownBranchType := loadField(snapshot, configdomain.KeyUnknownBranchType, configdomain.ParseBranchType)
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
