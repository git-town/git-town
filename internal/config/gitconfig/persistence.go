package gitconfig

import (
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Persistence provides high-level access to the Git CLI.
type Persistence struct {
	io IO
}

func (self Persistence) RemoteURL(remote gitdomain.Remote) Option[string] {
	output, err := self.io.Shell.Query("git", "remote", "get-url", remote.String())
	if err != nil {
		// NOTE: it's okay to ignore the error here.
		// If we get an error here, we simply don't use the origin remote.
		return None[string]()
	}
	return NewOption(strings.TrimSpace(output))
}

func (self Persistence) RemoveBranchTypeOverride(branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return self.io.RemoveConfigValue(configdomain.ConfigScopeLocal, key.Key)
}

func (self Persistence) RemoveCreatePrototypeBranches() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyDeprecatedCreatePrototypeBranches)
}

func (self Persistence) RemoveDevRemote() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyDevRemote)
}

func (self Persistence) RemoveFeatureRegex() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyFeatureRegex)
}

func (self Persistence) RemoveNewBranchType() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyNewBranchType)
}

func (self Persistence) RemoveParent(parent gitdomain.LocalBranchName) error {
	return self.io.RemoveLocalConfigValue(configdomain.NewParentKey(parent))
}

func (self Persistence) RemovePerennialBranches() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyPerennialBranches)
}

func (self Persistence) RemovePerennialRegex() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyPerennialRegex)
}

func (self Persistence) RemovePushHook() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyPushHook)
}

func (self Persistence) RemoveShareNewBranches() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyShareNewBranches)
}

func (self Persistence) RemoveShipDeleteTrackingBranch() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyShipDeleteTrackingBranch)
}

func (self Persistence) RemoveShipStrategy() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeyShipStrategy)
}

func (self Persistence) RemoveSyncFeatureStrategy() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeySyncFeatureStrategy)
}

func (self Persistence) RemoveSyncPerennialStrategy() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeySyncPerennialStrategy)
}

func (self Persistence) RemoveSyncPrototypeStrategy() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeySyncPrototypeStrategy)
}

func (self Persistence) RemoveSyncTags() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeySyncTags)
}

func (self Persistence) RemoveSyncUpstream() error {
	return self.io.RemoveLocalConfigValue(configdomain.KeySyncUpstream)
}

func (self Persistence) SetBranchTypeOverride(branch gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	return self.io.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String())
}

func (self Persistence) SetDevRemote(value gitdomain.Remote) error {
	return self.io.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, value.String())
}

func (self Persistence) SetFeatureRegex(value configdomain.FeatureRegex) error {
	return self.io.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, value.String())
}

func (self Persistence) SetNewBranchType(value configdomain.BranchType) error {
	return self.io.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType, value.String())
}

func (self Persistence) SetOffline(value configdomain.Offline) error {
	return self.io.SetConfigValue(configdomain.ConfigScopeGlobal, configdomain.KeyOffline, value.String())
}

func (self Persistence) SetParent(child, parent gitdomain.LocalBranchName) error {
	return self.io.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
}
