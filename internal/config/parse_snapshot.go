package config

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

// provides the branch type overrides stored in the given Git metadata snapshot
func NewBranchTypeOverridesInSnapshot(snapshot configdomain.SingleSnapshot, gitCommands gitconfig.IO) (configdomain.BranchTypeOverrides, error) {
	result := configdomain.BranchTypeOverrides{}
	for key, value := range snapshot.BranchTypeOverrideEntries() {
		branch := key.Branch()
		if branch == "" {
			// empty branch name --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigBranchTypeOverrideEmpty))
			_ = gitCommands.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			// empty branch type values are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigBranchTypeOverrideEmpty))
			_ = gitCommands.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
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

func NewLineageFromSnapshot(snapshot configdomain.SingleSnapshot, updateOutdated bool, gitCommands gitconfig.IO) (configdomain.Lineage, error) {
	result := configdomain.NewLineage()
	for key, value := range snapshot.LineageEntries() {
		child := key.ChildBranch()
		if child == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigLineageEmptyChild))
			_ = gitCommands.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigLineageEmptyChild))
			_ = gitCommands.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
			continue
		}
		if updateOutdated && child.String() == value {
			fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.ConfigLineageParentIsChild, child)))
			_ = gitCommands.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
		}
		parent := gitdomain.NewLocalBranchName(value)
		result = result.Set(child, parent)
	}
	return result, nil
}

func NewPartialConfigFromSnapshot(snapshot configdomain.SingleSnapshot, updateOutdated bool, gitCommands gitconfig.IO) (configdomain.PartialConfig, error) {
	aliases := snapshot.Aliases()
	branchTypeOverrides, err1 := NewBranchTypeOverridesInSnapshot(snapshot, gitCommands)
	contributionRegex, err2 := configdomain.ParseContributionRegex(snapshot[configdomain.KeyContributionRegex])
	featureRegex, err3 := configdomain.ParseFeatureRegex(snapshot[configdomain.KeyFeatureRegex])
	forgeType, err4 := forgedomain.ParseForgeType(snapshot[configdomain.KeyForgeType])
	githubConnectorType, err5 := forgedomain.ParseGitHubConnectorType(snapshot[configdomain.KeyGitHubConnectorType])
	gitlabConnectorType, err6 := forgedomain.ParseGitLabConnectorType(snapshot[configdomain.KeyGitLabConnectorType])
	lineage, err7 := NewLineageFromSnapshot(snapshot, updateOutdated, gitCommands)
	newBranchType, err8 := configdomain.ParseBranchType(snapshot[configdomain.KeyNewBranchType])
	observedRegex, err9 := configdomain.ParseObservedRegex(snapshot[configdomain.KeyObservedRegex])
	offline, err10 := configdomain.ParseOffline(snapshot[configdomain.KeyOffline], configdomain.KeyOffline)
	perennialRegex, err11 := configdomain.ParsePerennialRegex(snapshot[configdomain.KeyPerennialRegex])
	pushHook, err12 := configdomain.ParsePushHook(snapshot[configdomain.KeyPushHook], configdomain.KeyPushHook)
	shareNewBranches, err13 := configdomain.ParseShareNewBranches(snapshot[configdomain.KeyShareNewBranches], configdomain.KeyShareNewBranches)
	shipDeleteTrackingBranch, err14 := configdomain.ParseShipDeleteTrackingBranch(snapshot[configdomain.KeyShipDeleteTrackingBranch], configdomain.KeyShipDeleteTrackingBranch)
	shipStrategy, err15 := configdomain.ParseShipStrategy(snapshot[configdomain.KeyShipStrategy])
	syncFeatureStrategy, err16 := configdomain.ParseSyncFeatureStrategy(snapshot[configdomain.KeySyncFeatureStrategy])
	syncPerennialStrategy, err17 := configdomain.ParseSyncPerennialStrategy(snapshot[configdomain.KeySyncPerennialStrategy])
	syncPrototypeStrategy, err18 := configdomain.ParseSyncPrototypeStrategy(snapshot[configdomain.KeySyncPrototypeStrategy])
	syncTags, err19 := configdomain.ParseSyncTags(snapshot[configdomain.KeySyncTags], configdomain.KeySyncTags)
	syncUpstream, err20 := configdomain.ParseSyncUpstream(snapshot[configdomain.KeySyncUpstream], configdomain.KeySyncUpstream)
	unknownBranchType, err21 := configdomain.ParseBranchType(snapshot[configdomain.KeyUnknownBranchType])
	return configdomain.PartialConfig{
		Aliases:                  aliases,
		BitbucketAppPassword:     forgedomain.ParseBitbucketAppPassword(snapshot[configdomain.KeyBitbucketAppPassword]),
		BitbucketUsername:        forgedomain.ParseBitbucketUsername(snapshot[configdomain.KeyBitbucketUsername]),
		BranchTypeOverrides:      branchTypeOverrides,
		CodebergToken:            forgedomain.ParseCodebergToken(snapshot[configdomain.KeyCodebergToken]),
		ContributionRegex:        contributionRegex,
		DevRemote:                gitdomain.NewRemote(snapshot[configdomain.KeyDevRemote]),
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
	}, cmp.Or(err1, err2, err3, err4, err5, err6, err7, err8, err9, err10, err11, err12, err13, err14, err15, err16, err17, err18, err19, err20, err21)
}
