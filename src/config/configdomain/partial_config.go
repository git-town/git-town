package configdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                   Aliases
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

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{ //nolint:exhaustruct
		Aliases: Aliases{},
	}
}
