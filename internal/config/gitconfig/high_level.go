package gitconfig

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

func RemoveBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
}

func RemoveDevRemote(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote)
}

func RemoveFeatureRegex(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex)
}

func RemoveMainBranch(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyMainBranch)
}

func RemoveNewBranchType(runner subshelldomain.Runner) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType)
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

func SetBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String())
}

func SetDevRemote(runner subshelldomain.Runner, remote gitdomain.Remote) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, remote.String())
}

func SetFeatureRegex(runner subshelldomain.Runner, regex configdomain.FeatureRegex) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, regex.String())
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
