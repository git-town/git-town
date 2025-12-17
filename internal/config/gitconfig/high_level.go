package gitconfig

import (
	"strconv"
	"strings"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func DefaultBranch(querier subshelldomain.Querier) Option[gitdomain.LocalBranchName] {
	name, err := querier.QueryTrim("git", "config", "--get", "init.defaultbranch")
	if err != nil {
		return None[gitdomain.LocalBranchName]()
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return None[gitdomain.LocalBranchName]()
	}
	return Some(gitdomain.LocalBranchName(name))
}

func DefaultRemote(querier subshelldomain.Querier) gitdomain.Remote {
	output, err := querier.QueryTrim("git", "config", "--get", "clone.defaultRemoteName")
	if err != nil {
		// Git returns an error if the user has not configured a default remote name.
		// In this case use the Git default of "origin".
		return gitdomain.RemoteOrigin
	}
	return gitdomain.Remote(output)
}

func RemoveAlias(runner subshelldomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeGlobal, aliasableCommand.Key().Key())
}

func RemoveAutoSync(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyAutoSync)
}

func RemoveBitbucketAppPassword(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyBitbucketAppPassword)
}

func RemoveBitbucketUsername(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyBitbucketUsername)
}

func RemoveBranchPrefix(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyBranchPrefix)
}

func RemoveBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
}

func RemoveContributionRegex(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyContributionRegex)
}

func RemoveDetached(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDetached)
}

func RemoveDevRemote(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote)
}

func RemoveFeatureRegex(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex)
}

func RemoveForgeType(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyForgeType)
}

func RemoveForgejoToken(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyForgejoToken)
}

func RemoveGitHubConnectorType(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGitHubConnectorType)
}

func RemoveGitHubToken(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGitHubToken)
}

func RemoveGitLabConnectorType(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGitLabConnectorType)
}

func RemoveGitLabToken(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGitLabToken)
}

func RemoveGiteaToken(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGiteaToken)
}

func RemoveMainBranch(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyMainBranch)
}

func RemoveNewBranchType(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType)
}

func RemoveObservedRegex(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyObservedRegex)
}

func RemoveOrder(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyOrder)
}

func RemoveOriginHostname(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyHostingOriginHostname)
}

func RemoveParent(runner subshelldomain.Runner, child gitdomain.LocalBranchName) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child))
}

func RemovePerennialBranches(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches)
}

func RemovePerennialRegex(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialRegex)
}

func RemoveProposalsShowLineage(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyProposalsShowLineage)
}

func RemovePushBranches(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPushBranches)
}

func RemovePushHook(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPushHook)
}

func RemoveShareNewBranches(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShareNewBranches)
}

func RemoveShipDeleteTrackingBranch(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipDeleteTrackingBranch)
}

func RemoveShipIgnoreUncommitted(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipIgnoreUncommitted)
}

func RemoveShipStrategy(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipStrategy)
}

func RemoveStash(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyStash)
}

func RemoveSyncFeatureStrategy(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy)
}

func RemoveSyncPerennialStrategy(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy)
}

func RemoveSyncPrototypeStrategy(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy)
}

func RemoveSyncTags(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncTags)
}

func RemoveSyncUpstream(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncUpstream)
}

func RemoveUnknownBranchType(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyUnknownBranchType)
}

func SetAlias(runner subshelldomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return SetConfigValue(runner, configdomain.ConfigScopeGlobal, aliasableCommand.Key().Key(), "town "+aliasableCommand.String())
}

func SetAutoSync(runner subshelldomain.Runner, value configdomain.AutoSync, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyAutoSync, value.String())
}

func SetBitbucketAppPassword(runner subshelldomain.Runner, value forgedomain.BitbucketAppPassword, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyBitbucketAppPassword, value.String())
}

func SetBitbucketUsername(runner subshelldomain.Runner, value forgedomain.BitbucketUsername, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyBitbucketUsername, value.String())
}

func SetBranchPrefix(runner subshelldomain.Runner, value configdomain.BranchPrefix, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyBranchPrefix, value.String())
}

func SetBranchTypeOverride(runner subshelldomain.Runner, branchType configdomain.BranchType, branches ...gitdomain.LocalBranchName) error {
	for _, branch := range branches {
		if err := SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String()); err != nil {
			return err
		}
	}
	return nil
}

func SetContributionRegex(runner subshelldomain.Runner, regex configdomain.ContributionRegex, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyContributionRegex, regex.String())
}

func SetDetached(runner subshelldomain.Runner, value configdomain.Detached, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyDetached, value.String())
}

func SetDevRemote(runner subshelldomain.Runner, remote gitdomain.Remote, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyDevRemote, remote.String())
}

func SetFeatureRegex(runner subshelldomain.Runner, regex configdomain.FeatureRegex, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyFeatureRegex, regex.String())
}

func SetForgeType(runner subshelldomain.Runner, forgeType forgedomain.ForgeType, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyForgeType, forgeType.String())
}

func SetForgejoToken(runner subshelldomain.Runner, value forgedomain.ForgejoToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyForgejoToken, value.String())
}

func SetGitHubConnectorType(runner subshelldomain.Runner, value forgedomain.GitHubConnectorType, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGitHubConnectorType, value.String())
}

func SetGitHubToken(runner subshelldomain.Runner, value forgedomain.GitHubToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGitHubToken, value.String())
}

func SetGitLabConnectorType(runner subshelldomain.Runner, value forgedomain.GitLabConnectorType, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGitLabConnectorType, value.String())
}

func SetGitLabToken(runner subshelldomain.Runner, value forgedomain.GitLabToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGitLabToken, value.String())
}

func SetGiteaToken(runner subshelldomain.Runner, value forgedomain.GiteaToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGiteaToken, value.String())
}

func SetMainBranch(runner subshelldomain.Runner, value gitdomain.LocalBranchName, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyMainBranch, value.String())
}

func SetNewBranchType(runner subshelldomain.Runner, value configdomain.NewBranchType, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyNewBranchType, value.String())
}

func SetObservedRegex(runner subshelldomain.Runner, regex configdomain.ObservedRegex, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyObservedRegex, regex.String())
}

func SetOffline(runner subshelldomain.Runner, value configdomain.Offline) error {
	return SetConfigValue(runner, configdomain.ConfigScopeGlobal, configdomain.KeyOffline, value.String())
}

func SetOrder(runner subshelldomain.Runner, value configdomain.Order, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyOrder, value.String())
}

func SetOriginHostname(runner subshelldomain.Runner, hostname configdomain.HostingOriginHostname, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyHostingOriginHostname, hostname.String())
}

func SetParent(runner subshelldomain.Runner, child, parent gitdomain.LocalBranchName) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
}

func SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyPerennialBranches, branches.Join(" "))
}

func SetPerennialRegex(runner subshelldomain.Runner, value configdomain.PerennialRegex, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyPerennialRegex, value.String())
}

func SetProposalsShowLineage(runner subshelldomain.Runner, value forgedomain.ProposalsShowLineage, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyProposalsShowLineage, value.String())
}

func SetPushBranches(runner subshelldomain.Runner, value configdomain.PushBranches, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyPushBranches, strconv.FormatBool(bool(value)))
}

func SetPushHook(runner subshelldomain.Runner, value configdomain.PushHook, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

func SetShareNewBranches(runner subshelldomain.Runner, value configdomain.ShareNewBranches, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyShareNewBranches, value.String())
}

func SetShipDeleteTrackingBranch(runner subshelldomain.Runner, value configdomain.ShipDeleteTrackingBranch, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.ShouldDeleteTrackingBranch()))
}

func SetShipIgnoreUncommitted(runner subshelldomain.Runner, value configdomain.IgnoreUncommitted, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyShipIgnoreUncommitted, strconv.FormatBool(value.ShouldIgnoreUncommitted()))
}

func SetShipStrategy(runner subshelldomain.Runner, value configdomain.ShipStrategy, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyShipStrategy, value.String())
}

func SetStash(runner subshelldomain.Runner, value configdomain.Stash, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyStash, value.String())
}

func SetSyncFeatureStrategy(runner subshelldomain.Runner, value configdomain.SyncFeatureStrategy, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeySyncFeatureStrategy, value.String())
}

func SetSyncPerennialStrategy(runner subshelldomain.Runner, value configdomain.SyncPerennialStrategy, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeySyncPerennialStrategy, value.String())
}

func SetSyncPrototypeStrategy(runner subshelldomain.Runner, value configdomain.SyncPrototypeStrategy, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeySyncPrototypeStrategy, value.String())
}

func SetSyncTags(runner subshelldomain.Runner, value configdomain.SyncTags, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeySyncTags, value.String())
}

func SetSyncUpstream(runner subshelldomain.Runner, value configdomain.SyncUpstream, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeySyncUpstream, strconv.FormatBool(value.ShouldSyncUpstream()))
}

func SetUnknownBranchType(runner subshelldomain.Runner, value configdomain.UnknownBranchType, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyUnknownBranchType, value.String())
}
