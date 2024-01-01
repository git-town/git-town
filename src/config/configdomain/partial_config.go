package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                   map[AliasableCommand]string
	CodeHostingOriginHostname *CodeHostingOriginHostname
	CodeHostingPlatformName   *CodeHostingPlatformName
	GiteaToken                *GiteaToken
	GitHubToken               *GitHubToken
	GitLabToken               *GitLabToken
	GitUserEmail              *string
	GitUserName               *string
	Lineage                   *Lineage
	MainBranch                *gitdomain.LocalBranchName
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

func (self *PartialConfig) Add(key Key, value string) error {
	if strings.HasPrefix(key.name, "alias.") {
		aliasableCommand := LookupAliasableCommand(key)
		if aliasableCommand != nil {
			self.Aliases[*aliasableCommand] = value
		}
		return nil
	}
	if strings.HasPrefix(key.name, "git-town-branch.") {
		if self.Lineage == nil {
			self.Lineage = &Lineage{}
		}
		child := gitdomain.NewLocalBranchName(strings.TrimSuffix(strings.TrimPrefix(key.String(), "git-town-branch."), ".parent"))
		parent := gitdomain.NewLocalBranchName(value)
		(*self.Lineage)[child] = parent
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
	case KeyGitUserEmail:
		self.GitUserEmail = &value
	case KeyGitUserName:
		self.GitUserName = &value
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
		Aliases: map[AliasableCommand]string{},
	}
}
