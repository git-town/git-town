package config

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/pkg/colors"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// provides the branch type overrides stored in the given Git metadata snapshot
func NewBranchTypeOverridesInSnapshot(snapshot configdomain.SingleSnapshot, runner subshelldomain.Runner) (configdomain.BranchTypeOverrides, error) {
	result := configdomain.BranchTypeOverrides{}
	for key, value := range snapshot.BranchTypeOverrideEntries() {
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
	for key, value := range snapshot.LineageEntries() {
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

func NewPartialConfigFromSnapshot(snapshot configdomain.SingleSnapshot, updateOutdated bool, runner subshelldomain.Runner) (configdomain.PartialConfig, error) {
	aliases := snapshot.Aliases()
	autoResolve, errAutoResolve := configdomain.ParseAutoResolve(snapshot[configdomain.KeyAutoResolve], configdomain.KeyAutoResolve)
	branchTypeOverrides, errBranchTypeOverride := NewBranchTypeOverridesInSnapshot(snapshot, runner)
	contributionRegex, errContributionRegex := configdomain.ParseContributionRegex(snapshot[configdomain.KeyContributionRegex])
	featureRegex, errFeatureRegex := configdomain.ParseFeatureRegex(snapshot[configdomain.KeyFeatureRegex])
	forgeType, errForgeType := forgedomain.ParseForgeType(snapshot[configdomain.KeyForgeType])
	githubConnectorType, errGitHubConnectorType := forgedomain.ParseGitHubConnectorType(snapshot[configdomain.KeyGitHubConnectorType])
	gitlabConnectorType, errGitLabConnectorType := forgedomain.ParseGitLabConnectorType(snapshot[configdomain.KeyGitLabConnectorType])
	lineage, errLineage := NewLineageFromSnapshot(snapshot, updateOutdated, runner)
	newBranchTypeValue, errNewBranchType := configdomain.ParseBranchType(snapshot[configdomain.KeyNewBranchType])
	newBranchType := configdomain.NewBranchTypeOpt(newBranchTypeValue)
	observedRegex, errObservedRegex := configdomain.ParseObservedRegex(snapshot[configdomain.KeyObservedRegex])
	offline, errOffline := configdomain.ParseOffline(snapshot[configdomain.KeyOffline], configdomain.KeyOffline)
	perennialRegex, errPerennialRegex := configdomain.ParsePerennialRegex(snapshot[configdomain.KeyPerennialRegex])
	proposalsShowLineage, errProposalsShowLineage := forgedomain.ParseProposalsShowLineage(snapshot[configdomain.KeyProposalsShowLineage])
	pushHook, errPushHook := configdomain.ParsePushHook(snapshot[configdomain.KeyPushHook], configdomain.KeyPushHook)
	shareNewBranches, errShareNewBranches := configdomain.ParseShareNewBranches(snapshot[configdomain.KeyShareNewBranches], configdomain.KeyShareNewBranches)
	shipDeleteTrackingBranch, errShipDeleteTrackingBranch := configdomain.ParseShipDeleteTrackingBranch(snapshot[configdomain.KeyShipDeleteTrackingBranch], configdomain.KeyShipDeleteTrackingBranch)
	shipStrategy, errShipStrategy := configdomain.ParseShipStrategy(snapshot[configdomain.KeyShipStrategy])
	syncFeatureStrategy, errSyncFeatureStrategy := configdomain.ParseSyncFeatureStrategy(snapshot[configdomain.KeySyncFeatureStrategy])
	syncPerennialStrategy, errSyncPerennialStrategy := configdomain.ParseSyncPerennialStrategy(snapshot[configdomain.KeySyncPerennialStrategy])
	syncPrototypeStrategy, errSyncPrototypeStrategy := configdomain.ParseSyncPrototypeStrategy(snapshot[configdomain.KeySyncPrototypeStrategy])
	syncTags, errSyncTags := configdomain.ParseSyncTags(snapshot[configdomain.KeySyncTags], configdomain.KeySyncTags)
	syncUpstream, errSyncUpstream := configdomain.ParseSyncUpstream(snapshot[configdomain.KeySyncUpstream], configdomain.KeySyncUpstream)
	unknownBranchTypeValue, errUnknownBranchType := configdomain.ParseBranchType(snapshot[configdomain.KeyUnknownBranchType])
	unknownBranchType := configdomain.UnknownBranchTypeOpt(unknownBranchTypeValue)
	return configdomain.PartialConfig{
		Aliases:                  aliases,
		AutoResolve:              autoResolve,
		BitbucketAppPassword:     forgedomain.ParseBitbucketAppPassword(snapshot[configdomain.KeyBitbucketAppPassword]),
		BitbucketUsername:        forgedomain.ParseBitbucketUsername(snapshot[configdomain.KeyBitbucketUsername]),
		BranchTypeOverrides:      branchTypeOverrides,
		CodebergToken:            forgedomain.ParseCodebergToken(snapshot[configdomain.KeyCodebergToken]),
		ContributionRegex:        contributionRegex,
		DevRemote:                gitdomain.NewRemote(snapshot[configdomain.KeyDevRemote]),
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
		Offline:                  offline,
		PerennialBranches:        gitdomain.ParseLocalBranchNames(snapshot[configdomain.KeyPerennialBranches]),
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
		UnknownBranchType:        unknownBranchType,
		Verbose:                  None[configdomain.Verbose](),
	}, cmp.Or(errAutoResolve, errBranchTypeOverride, errContributionRegex, errFeatureRegex, errForgeType, errGitHubConnectorType, errGitLabConnectorType, errLineage, errNewBranchType, errObservedRegex, errOffline, errPerennialRegex, errProposalsShowLineage, errPushHook, errShareNewBranches, errShipDeleteTrackingBranch, errShipStrategy, errSyncFeatureStrategy, errSyncPerennialStrategy, errSyncPrototypeStrategy, errSyncTags, errSyncUpstream, errUnknownBranchType)
}
