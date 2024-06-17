package gitconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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
	KeyContributionBranches                = Key("git-town.contribution-branches")
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
	KeyGiteaToken                          = Key("git-town.gitea-token")
	KeyGithubToken                         = Key("git-town.github-token")
	KeyGitlabToken                         = Key("git-town.gitlab-token")
	KeyHostingOriginHostname               = Key("git-town.hosting-origin-hostname")
	KeyHostingPlatform                     = Key("git-town.hosting-platform")
	KeyMainBranch                          = Key("git-town.main-branch")
	KeyObservedBranches                    = Key("git-town.observed-branches")
	KeyOffline                             = Key("git-town.offline")
	KeyParkedBranches                      = Key("git-town.parked-branches")
	KeyPerennialBranches                   = Key("git-town.perennial-branches")
	KeyPerennialRegex                      = Key("git-town.perennial-regex")
	KeyPushHook                            = Key("git-town.push-hook")
	KeyPushNewBranches                     = Key("git-town.push-new-branches")
	KeyShipDeleteTrackingBranch            = Key("git-town.ship-delete-tracking-branch")
	KeySyncFeatureStrategy                 = Key("git-town.sync-feature-strategy")
	KeySyncPerennialStrategy               = Key("git-town.sync-perennial-strategy")
	KeySyncStrategy                        = Key("git-town.sync-strategy")
	KeySyncUpstream                        = Key("git-town.sync-upstream")
	KeyGitUserEmail                        = Key("user.email")
	KeyGitUserName                         = Key("user.name")
)

var keys = []Key{ //nolint:gochecknoglobals
	KeyHostingOriginHostname,
	KeyHostingPlatform,
	KeyContributionBranches,
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
	KeyGiteaToken,
	KeyGithubToken,
	KeyGitlabToken,
	KeyGitUserEmail,
	KeyGitUserName,
	KeyMainBranch,
	KeyObservedBranches,
	KeyOffline,
	KeyParkedBranches,
	KeyPerennialBranches,
	KeyPerennialRegex,
	KeyPushHook,
	KeyPushNewBranches,
	KeyShipDeleteTrackingBranch,
	KeySyncFeatureStrategy,
	KeySyncPerennialStrategy,
	KeySyncStrategy,
	KeySyncUpstream,
}

func AliasableCommandForKey(key Key) Option[configdomain.AliasableCommand] {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		if KeyForAliasableCommand(aliasableCommand) == key {
			return Some(aliasableCommand)
		}
	}
	return None[configdomain.AliasableCommand]()
}

func KeyForAliasableCommand(aliasableCommand configdomain.AliasableCommand) Key {
	switch aliasableCommand {
	case configdomain.AliasableCommandAppend:
		return KeyAliasAppend
	case configdomain.AliasableCommandCompress:
		return KeyAliasCompress
	case configdomain.AliasableCommandContribute:
		return KeyAliasContribute
	case configdomain.AliasableCommandDiffParent:
		return KeyAliasDiffParent
	case configdomain.AliasableCommandHack:
		return KeyAliasHack
	case configdomain.AliasableCommandKill:
		return KeyAliasKill
	case configdomain.AliasableCommandObserve:
		return KeyAliasObserve
	case configdomain.AliasableCommandPark:
		return KeyAliasPark
	case configdomain.AliasableCommandPrepend:
		return KeyAliasPrepend
	case configdomain.AliasableCommandPropose:
		return KeyAliasPropose
	case configdomain.AliasableCommandRenameBranch:
		return KeyAliasRenameBranch
	case configdomain.AliasableCommandRepo:
		return KeyAliasRepo
	case configdomain.AliasableCommandSetParent:
		return KeyAliasSetParent
	case configdomain.AliasableCommandShip:
		return KeyAliasShip
	case configdomain.AliasableCommandSync:
		return KeyAliasSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", &aliasableCommand))
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
	if lineageKey, hasLineageKey := parseLineageKey(name).Get(); hasLineageKey {
		return Some(lineageKey)
	}
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		key := KeyForAliasableCommand(aliasableCommand)
		if key.String() == name {
			return Some(key)
		}
	}
	return None[Key]()
}

const (
	LineageKeyPrefix = "git-town-branch."
	LineageKeySuffix = ".parent"
)

func parseLineageKey(key string) Option[Key] {
	if strings.HasPrefix(key, LineageKeyPrefix) && strings.HasSuffix(key, LineageKeySuffix) {
		result := Key(key)
		return Some(result)
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
