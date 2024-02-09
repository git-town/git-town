package gitconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
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
	KeyAliasDiffParent                     = Key("alias.diff-parent")
	KeyAliasHack                           = Key("alias.hack")
	KeyAliasKill                           = Key("alias.kill")
	KeyAliasPrepend                        = Key("alias.prepend")
	KeyAliasPropose                        = Key("alias.propose")
	KeyAliasRenameBranch                   = Key("alias.rename-branch")
	KeyAliasRepo                           = Key("alias.repo")
	KeyAliasSetParent                      = Key("alias.set-parent")
	KeyAliasShip                           = Key("alias.ship")
	KeyAliasSync                           = Key("alias.sync")
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
	KeyOffline                             = Key("git-town.offline")
	KeyPerennialBranches                   = Key("git-town.perennial-branches")
	KeyPushHook                            = Key("git-town.push-hook")
	KeyPushNewBranches                     = Key("git-town.push-new-branches")
	KeyShipDeleteTrackingBranch            = Key("git-town.ship-delete-tracking-branch")
	KeySyncBeforeShip                      = Key("git-town.sync-before-ship")
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
	KeyOffline,
	KeyPerennialBranches,
	KeyPushHook,
	KeyPushNewBranches,
	KeyShipDeleteTrackingBranch,
	KeySyncBeforeShip,
	KeySyncFeatureStrategy,
	KeySyncPerennialStrategy,
	KeySyncStrategy,
	KeySyncUpstream,
}

func AliasableCommandForKey(key Key) *configdomain.AliasableCommand {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		if KeyForAliasableCommand(aliasableCommand) == key {
			return &aliasableCommand
		}
	}
	return nil
}

func KeyForAliasableCommand(aliasableCommand configdomain.AliasableCommand) Key {
	switch aliasableCommand {
	case configdomain.AliasableCommandAppend:
		return KeyAliasAppend
	case configdomain.AliasableCommandDiffParent:
		return KeyAliasDiffParent
	case configdomain.AliasableCommandHack:
		return KeyAliasHack
	case configdomain.AliasableCommandKill:
		return KeyAliasKill
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
	return Key(fmt.Sprintf("git-town-branch.%s.parent", branch))
}

func ParseKey(name string) *Key {
	for _, configKey := range keys {
		if configKey.String() == name {
			return &configKey
		}
	}
	lineageKey := parseLineageKey(name)
	if lineageKey != nil {
		return lineageKey
	}
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		key := KeyForAliasableCommand(aliasableCommand)
		if key.String() == name {
			return &key
		}
	}
	return nil
}

func parseLineageKey(key string) *Key {
	if !strings.HasPrefix(key, "git-town-branch.") || !strings.HasSuffix(key, ".parent") {
		return nil
	}
	result := Key(key)
	return &result
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
