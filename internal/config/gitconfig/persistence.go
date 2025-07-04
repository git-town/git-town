package gitconfig

import (
	"strconv"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Persistence provides high-level access to the Git CLI.
type Persistence struct {
	IO IO
}

func (self Persistence) RemoteURL(remote gitdomain.Remote) Option[string] {
	output, err := self.IO.Shell.Query("git", "remote", "get-url", remote.String())
	if err != nil {
		// NOTE: it's okay to ignore the error here.
		// If we get an error here, we simply don't use the origin remote.
		return None[string]()
	}
	return NewOption(strings.TrimSpace(output))
}

func (self Persistence) RemoveBranchTypeOverride(branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return self.IO.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
}

func (self Persistence) RemoveCreatePrototypeBranches() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyDeprecatedCreatePrototypeBranches)
}

func (self Persistence) RemoveDevRemote() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyDevRemote)
}

func (self Persistence) RemoveFeatureRegex() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyFeatureRegex)
}

func (self Persistence) RemoveForgeType() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyForgeType)
}

func (self Persistence) RemoveMainBranch() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyMainBranch)
}

func (self Persistence) RemoveNewBranchType() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyNewBranchType)
}

func (self Persistence) RemoveParent(parent gitdomain.LocalBranchName) error {
	return self.IO.RemoveLocalConfigValue(configdomain.NewParentKey(parent))
}

func (self Persistence) RemovePerennialBranches() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyPerennialBranches)
}

func (self Persistence) RemovePerennialRegex() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyPerennialRegex)
}

func (self Persistence) RemovePushHook() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyPushHook)
}

func (self Persistence) RemoveShareNewBranches() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyShareNewBranches)
}

func (self Persistence) RemoveShipDeleteTrackingBranch() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyShipDeleteTrackingBranch)
}

func (self Persistence) RemoveShipStrategy() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeyShipStrategy)
}

func (self Persistence) RemoveSyncFeatureStrategy() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeySyncFeatureStrategy)
}

func (self Persistence) RemoveSyncPerennialStrategy() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeySyncPerennialStrategy)
}

func (self Persistence) RemoveSyncPrototypeStrategy() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeySyncPrototypeStrategy)
}

func (self Persistence) RemoveSyncTags() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeySyncTags)
}

func (self Persistence) RemoveSyncUpstream() error {
	return self.IO.RemoveLocalConfigValue(configdomain.KeySyncUpstream)
}

func (self Persistence) SetBranchTypeOverride(branch gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String())
}

func (self Persistence) SetDevRemote(value gitdomain.Remote) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, value.String())
}

func (self Persistence) SetFeatureRegex(value configdomain.FeatureRegex) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, value.String())
}

func (self Persistence) SetForgeType(value forgedomain.ForgeType) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyForgeType, value.String())
}

func (self Persistence) SetMainBranch(value gitdomain.LocalBranchName) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyMainBranch, value.String())
}

func (self Persistence) SetNewBranchType(value configdomain.BranchType) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType, value.String())
}

func (self Persistence) SetOffline(value configdomain.Offline) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeGlobal, configdomain.KeyOffline, value.String())
}

func (self Persistence) SetParent(child, parent gitdomain.LocalBranchName) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
}

func (self Persistence) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches, branches.Join(" "))
}

func (self Persistence) SetPerennialRegex(value configdomain.PerennialRegex) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPerennialRegex, value.String())
}

func (self Persistence) SetPushHook(value configdomain.PushHook) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPushHook, strconv.FormatBool(value.IsTrue()))
}

func (self Persistence) SetShareNewBranches(value configdomain.ShareNewBranches, scope configdomain.ConfigScope) error {
	return self.IO.SetConfigValue(scope, configdomain.KeyShareNewBranches, value.String())
}

func (self Persistence) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, scope configdomain.ConfigScope) error {
	return self.IO.SetConfigValue(scope, configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.IsTrue()))
}

func (self Persistence) SetShipStrategy(value configdomain.ShipStrategy, scope configdomain.ConfigScope) error {
	return self.IO.SetConfigValue(scope, configdomain.KeyShipStrategy, value.String())
}

func (self Persistence) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy, value.String())
}

func (self Persistence) SetSyncPerennialStrategy(value configdomain.SyncPerennialStrategy) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy, value.String())
}

func (self Persistence) SetSyncPrototypeStrategy(value configdomain.SyncPrototypeStrategy) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy, value.String())
}

func (self Persistence) SetSyncTags(value configdomain.SyncTags) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncTags, value.String())
}

func (self Persistence) SetSyncUpstream(value configdomain.SyncUpstream, scope configdomain.ConfigScope) error {
	return self.IO.SetConfigValue(scope, configdomain.KeySyncUpstream, strconv.FormatBool(value.IsTrue()))
}

func (self Persistence) SetUnknownBranchType(value configdomain.BranchType) error {
	return self.IO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyUnknownBranchType, value.String())
}
