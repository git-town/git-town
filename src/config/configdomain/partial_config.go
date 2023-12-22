package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                   map[Key]string
	CodeHostingOriginHostname *CodeHostingOriginHostname
	CodeHostingPlatformName   *CodeHostingPlatformName
	GiteaToken                *GiteaToken
	GitHubToken               *GitHubToken
	GitLabToken               *GitLabToken
	Lineage                   Lineage
	MainBranch                *domain.LocalBranchName
	NewBranchPush             *NewBranchPush
	Offline                   *Offline
	PerennialBranches         *gitdomain.LocalBranchNames
	PushHook                  *PushHook
	ShipDeleteTrackingBranch  *ShipDeleteTrackingBranch
	SyncBeforeShip            *SyncBeforeShip
	SyncFeatureStrategy       *SyncFeatureStrategy
	SyncPerennialStrategy     *SyncPerennialStrategy
	SyncUpstream              *SyncUpstream
}

func (self *PartialConfig) Add(key Key, value string, deleteEntry func(Key) error) error {
	if strings.HasPrefix(key.name, "alias.") {
		self.Aliases[key] = value
		return nil
	}
	if strings.HasPrefix(key.name, "git-town-branch.") {
		child := domain.NewLocalBranchName(strings.TrimSuffix(strings.TrimPrefix(key.String(), "git-town-branch."), ".parent"))
		if value == "" {
			_ = deleteEntry(key)
			fmt.Printf("\nNOTICE: I have found an empty parent configuration entry for branch %q.\n", child)
			fmt.Println("I have deleted this configuration entry.")
		} else {
			self.Lineage[child] = domain.NewLocalBranchName(value)
		}
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
		self.MainBranch = gitdomain.NewLocalBranchNameRefAllowEmpty(value)
	case KeyOffline:
		self.Offline, err = NewOfflineRef(value)
	case KeyPerennialBranches:
		self.PerennialBranches = gitdomain.ParseLocalBranchNamesRef(value)
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
	return PartialConfig{ //nolint:exhaustruct
		Aliases: map[Key]string{},
		Lineage: Lineage{},
	}
}
