package gitconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

func RemoveBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
}

func RemoveDeprecatedCreatePrototypeBranches(runner subshelldomain.Runner) {
	_ = RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDeprecatedCreatePrototypeBranches)
}

func RemoveDevRemote(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyDevRemote)
}

func RemoveFeatureRegex(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyFeatureRegex)
}

func RemoveNewBranchType(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyNewBranchType)
}

func RemoveParent(runner subshelldomain.Runner, child gitdomain.LocalBranchName) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child))
}

func RemovePerennialBranches(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyPerennialBranches)
}

func RemovePerennialRegex(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyPerennialRegex)
}

func RemovePushHook(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyPushHook)
}

func RemoveShareNewBranches(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyShareNewBranches)
}

func RemoveShipDeleteTrackingBranch(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyShipDeleteTrackingBranch)
}

func RemoveShipStrategy(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyShipStrategy)
}

func RemoveSyncFeatureStrategy(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeySyncFeatureStrategy)
}

func RemoveSyncPerennialStrategy(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeySyncPerennialStrategy)
}

func RemoveSyncPrototypeStrategy(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeySyncPrototypeStrategy)
}

func RemoveSyncTags(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeySyncTags)
}

func RemoveSyncUpstream(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeySyncUpstream)
}

func SetBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String())
}

func SetDevRemote(runner subshelldomain.Runner, remote gitdomain.Remote) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, remote.String())
}

func SetFeatureRegex(runner subshelldomain.Runner, regex configdomain.FeatureRegex) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, regex.String())
}

func SetNewBranchType(runner subshelldomain.Runner, value configdomain.BranchType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType, value.String())
}

func SetOffline(runner subshelldomain.Runner, scope configdomain.ConfigScope, value configdomain.Offline) error {
	return SetConfigValue(runner, scope, configdomain.KeyOffline, value.String())
}

func SetParent(runner subshelldomain.Runner, child, parent gitdomain.LocalBranchName) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
}

func SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches, branches.Join(" "))
}
