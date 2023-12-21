package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v11/src/domain"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                   map[Key]string
	CodeHostingOriginHostname *CodeHostingOriginHostname
	CodeHostingPlatformName   *CodeHostingPlatformName
	GiteaToken                *GiteaToken
	GitHubToken               *GitHubToken
	GitLabToken               *GitLabToken
	MainBranch                *domain.LocalBranchName
	NewBranchPush             *NewBranchPush
	Offline                   *Offline
	PerennialBranches         *domain.LocalBranchNames
	PushHook                  *PushHook
	ShipDeleteTrackingBranch  *ShipDeleteTrackingBranch
	SyncBeforeShip            *SyncBeforeShip
	SyncFeatureStrategy       *SyncFeatureStrategy
	SyncPerennialStrategy     *SyncPerennialStrategy
	SyncUpstream              *SyncUpstream
}

func (self *PartialConfig) Add(key Key, value string) error {
	if strings.HasPrefix(key.name, "alias.") {
		self.Aliases[key] = value
		return nil
	}
	var err error
	switch key {
	case KeyCodeHostingOriginHostname:
		self.CodeHostingOriginHostname = NewCodeHostingOriginHostnameRef(value)
	case KeyCodeHostingPlatform:
		self.CodeHostingPlatformName = NewCodeHostingPlatformNameRef(value)
	case KeyGiteaToken:
		self.GiteaToken = NewGiteaTokenRef(value)
	case KeyGithubToken:
		self.GitHubToken = NewGitHubTokenRef(value)
	case KeyGitlabToken:
		self.GitLabToken = NewGitLabTokenRef(value)
	case KeyMainBranch:
		self.MainBranch = domain.NewLocalBranchNameRefAllowEmpty(value)
	case KeyOffline:
		self.Offline, err = NewOfflineRef(value)
	case KeyPerennialBranches:
		self.PerennialBranches = domain.ParseLocalBranchNamesRef(value)
	case KeyPushHook:
		self.PushHook, err = NewPushHookRef(value)
	case KeyPushNewBranches:
		self.NewBranchPush, err = ParseNewBranchPushRef(value)
	case KeyShipDeleteTrackingBranch:
		self.ShipDeleteTrackingBranch, err = ParseShipDeleteTrackingBranchRef(value)
	case KeySyncBeforeShip:
		self.SyncBeforeShip, err = NewSyncBeforeShipRef(value)
	case KeySyncFeatureStrategy:
		self.SyncFeatureStrategy, err = NewSyncFeatureStrategyRef(value)
	case KeySyncPerennialStrategy:
		self.SyncPerennialStrategy, err = NewSyncPerennialStrategyRef(value)
	case KeySyncUpstream:
		self.SyncUpstream, err = ParseSyncUpstreamRef(value)
	default:
		panic("unprocessed key: " + key.String())
	}
	return err
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: map[Key]string{},
	} //nolint:exhaustruct
}

// PartialConfigDiff diffs the given PartialConfig instances.
func PartialConfigDiff(before, after PartialConfig) ConfigDiff {
	result := ConfigDiff{
		Added:   []Key{},
		Removed: map[Key]string{},
		Changed: map[Key]domain.Change[string]{},
	}
	DiffPtr(&result, KeyCodeHostingOriginHostname, before.CodeHostingOriginHostname, after.CodeHostingOriginHostname)
	DiffPtr(&result, KeyCodeHostingPlatform, before.CodeHostingPlatformName, after.CodeHostingPlatformName)
	DiffPtr(&result, KeyGiteaToken, before.GiteaToken, after.GiteaToken)
	DiffPtr(&result, KeyGithubToken, before.GitHubToken, after.GitHubToken)
	DiffPtr(&result, KeyGitlabToken, before.GitLabToken, after.GitLabToken)
	DiffPtr(&result, KeyMainBranch, before.MainBranch, after.MainBranch)
	DiffPtr(&result, KeyOffline, before.Offline, after.Offline)
	DiffLocalBranchNames(&result, KeyPerennialBranches, before.PerennialBranches, after.PerennialBranches)
	DiffPtr(&result, KeyPushHook, before.PushHook, after.PushHook)
	DiffPtr(&result, KeyPushNewBranches, before.NewBranchPush, after.NewBranchPush)
	DiffPtr(&result, KeyShipDeleteTrackingBranch, before.ShipDeleteTrackingBranch, after.ShipDeleteTrackingBranch)
	DiffPtr(&result, KeySyncFeatureStrategy, before.SyncFeatureStrategy, after.SyncFeatureStrategy)
	DiffPtr(&result, KeySyncPerennialStrategy, before.SyncPerennialStrategy, after.SyncPerennialStrategy)
	DiffPtr(&result, KeySyncUpstream, before.SyncUpstream, after.SyncUpstream)
	return result
}
