package gitconfig

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// TODO: make this a method of SingleSnapshot?
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

// TODO: make this a method of SingleSnapshot?
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

func RemoveBitbucketAppPassword(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyBitbucketAppPassword)
}

func RemoveBitbucketUsername(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyBitbucketUsername)
}

func RemoveBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
}

func RemoveCodebergToken(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyCodebergToken)
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

func RemovePushHook(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPushHook)
}

func RemoveShareNewBranches(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShareNewBranches)
}

func RemoveShipDeleteTrackingBranch(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipDeleteTrackingBranch)
}

func RemoveShipStrategy(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipStrategy)
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

func SetAlias(runner subshelldomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return SetConfigValue(runner, configdomain.ConfigScopeGlobal, aliasableCommand.Key().Key(), "town "+aliasableCommand.String())
}

func SetBitbucketAppPassword(runner subshelldomain.Runner, value forgedomain.BitbucketAppPassword, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyBitbucketAppPassword, value.String())
}

func SetBitbucketUsername(runner subshelldomain.Runner, value forgedomain.BitbucketUsername, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyBitbucketUsername, value.String())
}

func SetBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String())
}

func SetCodebergToken(runner subshelldomain.Runner, value forgedomain.CodebergToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyCodebergToken, value.String())
}

func SetDevRemote(runner subshelldomain.Runner, remote gitdomain.Remote) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, remote.String())
}

func SetFeatureRegex(runner subshelldomain.Runner, regex configdomain.FeatureRegex) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, regex.String())
}

func SetForgeType(runner subshelldomain.Runner, forgeType forgedomain.ForgeType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyForgeType, forgeType.String())
}

func SetGitHubConnectorType(runner subshelldomain.Runner, value forgedomain.GitHubConnectorType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGitHubConnectorType, value.String())
}

func SetGitHubToken(runner subshelldomain.Runner, value forgedomain.GitHubToken, scope configdomain.ConfigScope) error {
	fmt.Println("111111111111111111111111111111", scope)
	return SetConfigValue(runner, scope, configdomain.KeyGitHubToken, value.String())
}

func SetGitLabConnectorType(runner subshelldomain.Runner, value forgedomain.GitLabConnectorType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyGitLabConnectorType, value.String())
}

func SetGitLabToken(runner subshelldomain.Runner, value forgedomain.GitLabToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGitLabToken, value.String())
}

func SetGiteaToken(runner subshelldomain.Runner, value forgedomain.GiteaToken, scope configdomain.ConfigScope) error {
	return SetConfigValue(runner, scope, configdomain.KeyGiteaToken, value.String())
}

func SetMainBranch(runner subshelldomain.Runner, value gitdomain.LocalBranchName) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyMainBranch, value.String())
}

func SetNewBranchType(runner subshelldomain.Runner, value configdomain.BranchType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType, value.String())
}

func SetOffline(runner subshelldomain.Runner, value configdomain.Offline) error {
	return SetConfigValue(runner, configdomain.ConfigScopeGlobal, configdomain.KeyOffline, value.String())
}

func SetOriginHostname(runner subshelldomain.Runner, hostname configdomain.HostingOriginHostname) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyHostingOriginHostname, hostname.String())
}

func SetParent(runner subshelldomain.Runner, child, parent gitdomain.LocalBranchName) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
}

func SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches, branches.Join(" "))
}

func SetPerennialRegex(runner subshelldomain.Runner, value configdomain.PerennialRegex) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialRegex, value.String())
}

func SetPushHook(runner subshelldomain.Runner, value configdomain.PushHook) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

func SetShareNewBranches(runner subshelldomain.Runner, value configdomain.ShareNewBranches) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShareNewBranches, value.String())
}

func SetShipDeleteTrackingBranch(runner subshelldomain.Runner, value configdomain.ShipDeleteTrackingBranch) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.IsTrue()))
}

func SetShipStrategy(runner subshelldomain.Runner, value configdomain.ShipStrategy) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipStrategy, value.String())
}

func SetSyncFeatureStrategy(runner subshelldomain.Runner, value configdomain.SyncFeatureStrategy) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy, value.String())
}

func SetSyncPerennialStrategy(runner subshelldomain.Runner, value configdomain.SyncPerennialStrategy) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy, value.String())
}

func SetSyncPrototypeStrategy(runner subshelldomain.Runner, value configdomain.SyncPrototypeStrategy) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy, value.String())
}

func SetSyncTags(runner subshelldomain.Runner, value configdomain.SyncTags) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncTags, value.String())
}

func SetSyncUpstream(runner subshelldomain.Runner, value configdomain.SyncUpstream) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncUpstream, strconv.FormatBool(value.IsTrue()))
}

func SetUnknownBranchType(runner subshelldomain.Runner, value configdomain.BranchType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyUnknownBranchType, value.String())
}
