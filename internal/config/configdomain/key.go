package configdomain

import (
	"encoding/json"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/pkg"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Key contains all the keys used in Git Town's Git metadata configuration.
type Key string

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (self Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.String())
}

func (self Key) String() string { return string(self) }

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (self *Key) UnmarshalJSON(b []byte) error {
	value := ""
	err := json.Unmarshal(b, &value)
	*self = Key(value)
	return err
}

const (
	KeyAliasAppend                         = Key("alias.append")
	KeyAliasCompress                       = Key("alias.compress")
	KeyAliasContribute                     = Key("alias.contribute")
	KeyAliasDiffParent                     = Key("alias.diff-parent")
	KeyAliasHack                           = Key("alias.hack")
	KeyAliasKill                           = Key("alias.kill")
	KeyAliasObserve                        = Key("alias.observe")
	KeyAliasPark                           = Key("alias.park")
	KeyAliasPrepend                        = Key("alias.prepend")
	KeyAliasPropose                        = Key("alias.propose")
	KeyAliasRenameBranch                   = Key("alias.rename-branch")
	KeyAliasRepo                           = Key("alias.repo")
	KeyAliasSetParent                      = Key("alias.set-parent")
	KeyAliasShip                           = Key("alias.ship")
	KeyAliasSync                           = Key("alias.sync")
	KeyBitbucketAppPassword                = Key("git-town.bitbucket-app-password")
	KeyBitbucketUsername                   = Key("git-town.bitbucket-username")
	KeyContributionBranches                = Key("git-town.contribution-branches")
	KeyContributionRegex                   = Key("git-town.contribution-regex")
	KeyCreatePrototypeBranches             = Key("git-town.create-prototype-branches")
	KeyDefaultBranchType                   = Key("git-town.default-branch-type")
	KeyDeprecatedCodeHostingDriver         = Key("git-town.code-hosting-driver")
	KeyDeprecatedCodeHostingOriginHostname = Key("git-town.code-hosting-origin-hostname")
	KeyDeprecatedCodeHostingPlatform       = Key("git-town.code-hosting-platform")
	KeyDeprecatedMainBranchName            = Key("git-town.main-branch-name")
	KeyDeprecatedNewBranchPushFlag         = Key("git-town.new-branch-push-flag")
	KeyDeprecatedPerennialBranchNames      = Key("git-town.perennial-branch-names")
	KeyDeprecatedPullBranchStrategy        = Key("git-town.pull-branch-strategy")
	KeyDeprecatedPushVerify                = Key("git-town.push-verify")
	KeyDeprecatedShipDeleteRemoteBranch    = Key("git-town.ship-delete-remote-branch")
	KeyDeprecatedSyncStrategy              = Key("git-town.sync-strategy")
	KeyFeatureRegex                        = Key("git-town.feature-regex")
	KeyGiteaToken                          = Key("git-town.gitea-token")
	KeyGithubToken                         = Key(pkg.KeyGithubToken)
	KeyGitlabToken                         = Key("git-town.gitlab-token")
	KeyHostingOriginHostname               = Key("git-town.hosting-origin-hostname")
	KeyHostingPlatform                     = Key("git-town.hosting-platform")
	KeyMainBranch                          = Key("git-town.main-branch")
	KeyObservedBranches                    = Key("git-town.observed-branches")
	KeyObservedRegex                       = Key("git-town.observed-regex")
	KeyOffline                             = Key("git-town.offline")
	KeyParkedBranches                      = Key("git-town.parked-branches")
	KeyPerennialBranches                   = Key("git-town.perennial-branches")
	KeyPerennialRegex                      = Key("git-town.perennial-regex")
	KeyPrototypeBranches                   = Key("git-town.prototype-branches")
	KeyPushHook                            = Key("git-town.push-hook")
	KeyPushNewBranches                     = Key("git-town.push-new-branches")
	KeyShipDeleteTrackingBranch            = Key("git-town.ship-delete-tracking-branch")
	KeyShipStrategy                        = Key("git-town.ship-strategy")
	KeyObsoleteSyncBeforeShip              = Key("git-town.sync-before-ship")
	KeySyncFeatureStrategy                 = Key("git-town.sync-feature-strategy")
	KeySyncPerennialStrategy               = Key("git-town.sync-perennial-strategy")
	KeySyncPrototypeStrategy               = Key("git-town.sync-prototype-strategy")
	KeySyncStrategy                        = Key("git-town.sync-strategy")
	KeySyncTags                            = Key("git-town.sync-tags")
	KeySyncUpstream                        = Key("git-town.sync-upstream")
	KeyGitUserEmail                        = Key("user.email")
	KeyGitUserName                         = Key("user.name")
)

var keys = []Key{ //nolint:gochecknoglobals
	KeyHostingOriginHostname,
	KeyHostingPlatform,
	KeyBitbucketAppPassword,
	KeyBitbucketUsername,
	KeyContributionBranches,
	KeyContributionRegex,
	KeyCreatePrototypeBranches,
	KeyDefaultBranchType,
	KeyDeprecatedCodeHostingDriver,
	KeyDeprecatedCodeHostingOriginHostname,
	KeyDeprecatedCodeHostingPlatform,
	KeyDeprecatedMainBranchName,
	KeyDeprecatedNewBranchPushFlag,
	KeyDeprecatedPerennialBranchNames,
	KeyDeprecatedPullBranchStrategy,
	KeyDeprecatedPushVerify,
	KeyDeprecatedShipDeleteRemoteBranch,
	KeyDeprecatedSyncStrategy,
	KeyFeatureRegex,
	KeyGiteaToken,
	KeyGithubToken,
	KeyGitlabToken,
	KeyGitUserEmail,
	KeyGitUserName,
	KeyMainBranch,
	KeyObservedBranches,
	KeyObservedRegex,
	KeyOffline,
	KeyParkedBranches,
	KeyPerennialBranches,
	KeyPerennialRegex,
	KeyPrototypeBranches,
	KeyPushHook,
	KeyPushNewBranches,
	KeyShipDeleteTrackingBranch,
	KeyShipStrategy,
	KeyObsoleteSyncBeforeShip,
	KeySyncFeatureStrategy,
	KeySyncPerennialStrategy,
	KeySyncPrototypeStrategy,
	KeySyncStrategy,
	KeySyncTags,
	KeySyncUpstream,
}

func NewParentKey(branch gitdomain.LocalBranchName) Key {
	return Key(LineageKeyPrefix + branch + LineageKeySuffix)
}

func ParseKey(name string) Option[Key] {
	for _, configKey := range keys {
		if configKey.String() == name {
			return Some(configKey)
		}
	}
	if isLineageKey(name) {
		return Some(Key(name))
	}
	if aliasKey, isAliasKey := AllAliasableCommands().LookupKey(name).Get(); isAliasKey {
		return Some(aliasKey.Key())
	}
	return None[Key]()
}

// DeprecatedKeys defines the up-to-date counterparts to deprecated configuration settings.
var DeprecatedKeys = map[Key]Key{ //nolint:gochecknoglobals
	KeyDeprecatedCodeHostingDriver:         KeyHostingPlatform,
	KeyDeprecatedCodeHostingOriginHostname: KeyHostingOriginHostname,
	KeyDeprecatedCodeHostingPlatform:       KeyHostingPlatform,
	KeyDeprecatedMainBranchName:            KeyMainBranch,
	KeyDeprecatedNewBranchPushFlag:         KeyPushNewBranches,
	KeyDeprecatedPerennialBranchNames:      KeyPerennialBranches,
	KeyDeprecatedPullBranchStrategy:        KeySyncPerennialStrategy,
	KeyDeprecatedPushVerify:                KeyPushHook,
	KeyDeprecatedShipDeleteRemoteBranch:    KeyShipDeleteTrackingBranch,
	KeyDeprecatedSyncStrategy:              KeySyncFeatureStrategy,
}

// ObsoleteKeys defines the keys that are sunset and should get deleted
var ObsoleteKeys = []Key{ //nolint:gochecknoglobals
	KeyObsoleteSyncBeforeShip,
}
