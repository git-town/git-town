package gitconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// Key contains all the keys used in Git Town configuration.
type Key struct {
	name string
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (self Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.name)
}

func (self Key) String() string { return self.name }

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (self *Key) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &self.name)
}

var (
	KeyAliasAppend                    = Key{"alias.append"}                          //nolint:gochecknoglobals
	KeyAliasDiffParent                = Key{"alias.diff-parent"}                     //nolint:gochecknoglobals
	KeyAliasHack                      = Key{"alias.hack"}                            //nolint:gochecknoglobals
	KeyAliasKill                      = Key{"alias.kill"}                            //nolint:gochecknoglobals
	KeyAliasPrepend                   = Key{"alias.prepend"}                         //nolint:gochecknoglobals
	KeyAliasPropose                   = Key{"alias.propose"}                         //nolint:gochecknoglobals
	KeyAliasRenameBranch              = Key{"alias.rename-branch"}                   //nolint:gochecknoglobals
	KeyAliasRepo                      = Key{"alias.repo"}                            //nolint:gochecknoglobals
	KeyAliasShip                      = Key{"alias.ship"}                            //nolint:gochecknoglobals
	KeyAliasSync                      = Key{"alias.sync"}                            //nolint:gochecknoglobals
	KeyCodeHostingOriginHostname      = Key{"git-town.code-hosting-origin-hostname"} //nolint:gochecknoglobals
	KeyCodeHostingPlatform            = Key{"git-town.code-hosting-platform"}        //nolint:gochecknoglobals
	KeyDeprecatedCodeHostingDriver    = Key{"git-town.code-hosting-driver"}          //nolint:gochecknoglobals
	KeyDeprecatedMainBranchName       = Key{"git-town.main-branch-name"}             //nolint:gochecknoglobals
	KeyDeprecatedNewBranchPushFlag    = Key{"git-town.new-branch-push-flag"}         //nolint:gochecknoglobals
	KeyDeprecatedPerennialBranchNames = Key{"git-town.perennial-branch-names"}       //nolint:gochecknoglobals
	KeyDeprecatedPullBranchStrategy   = Key{"git-town.pull-branch-strategy"}         //nolint:gochecknoglobals
	KeyDeprecatedPushVerify           = Key{"git-town.push-verify"}                  //nolint:gochecknoglobals
	KeyDeprecatedSyncStrategy         = Key{"git-town.sync-strategy"}                //nolint:gochecknoglobals
	KeyGiteaToken                     = Key{"git-town.gitea-token"}                  //nolint:gochecknoglobals
	KeyGithubToken                    = Key{"git-town.github-token"}                 //nolint:gochecknoglobals
	KeyGitlabToken                    = Key{"git-town.gitlab-token"}                 //nolint:gochecknoglobals
	KeyMainBranch                     = Key{"git-town.main-branch"}                  //nolint:gochecknoglobals
	KeyOffline                        = Key{"git-town.offline"}                      //nolint:gochecknoglobals
	KeyPerennialBranches              = Key{"git-town.perennial-branches"}           //nolint:gochecknoglobals
	KeyPushHook                       = Key{"git-town.push-hook"}                    //nolint:gochecknoglobals
	KeyPushNewBranches                = Key{"git-town.push-new-branches"}            //nolint:gochecknoglobals
	KeyShipDeleteTrackingBranch       = Key{"git-town.ship-delete-remote-branch"}    //nolint:gochecknoglobals
	KeySyncBeforeShip                 = Key{"git-town.sync-before-ship"}             //nolint:gochecknoglobals
	KeySyncFeatureStrategy            = Key{"git-town.sync-feature-strategy"}        //nolint:gochecknoglobals
	KeySyncPerennialStrategy          = Key{"git-town.sync-perennial-strategy"}      //nolint:gochecknoglobals
	KeySyncStrategy                   = Key{"git-town.sync-strategy"}                //nolint:gochecknoglobals
	KeySyncUpstream                   = Key{"git-town.sync-upstream"}                //nolint:gochecknoglobals
	KeyTestingRemoteURL               = Key{"git-town.testing.remote-url"}           //nolint:gochecknoglobals
	KeyGitUserEmail                   = Key{"user.email"}                            //nolint:gochecknoglobals
	KeyGitUserName                    = Key{"user.name"}                             //nolint:gochecknoglobals
)

var keys = []Key{ //nolint:gochecknoglobals
	KeyCodeHostingOriginHostname,
	KeyCodeHostingPlatform,
	KeyDeprecatedCodeHostingDriver,
	KeyDeprecatedMainBranchName,
	KeyDeprecatedNewBranchPushFlag,
	KeyDeprecatedPerennialBranchNames,
	KeyDeprecatedPullBranchStrategy,
	KeyDeprecatedPushVerify,
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
	KeyTestingRemoteURL,
}

func AliasableCommandToKey(aliasableCommand configdomain.AliasableCommand) Key {
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
	case configdomain.AliasableCommandShip:
		return KeyAliasShip
	case configdomain.AliasableCommandSync:
		return KeyAliasSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", &aliasableCommand))
}

func KeyToAliasableCommand(key Key) *configdomain.AliasableCommand {
	for _, aliasableCommand := range configdomain.AliasableCommands() {
		if AliasableCommandToKey(aliasableCommand) == key {
			return &aliasableCommand
		}
	}
	return nil
}

func NewKey(name string) Key {
	return Key{name}
}

func NewParentKey(branch gitdomain.LocalBranchName) Key {
	return Key{
		name: fmt.Sprintf("git-town-branch.%s.parent", branch),
	}
}

func ParseKey(name string) *Key {
	for _, configKey := range keys {
		if configKey.name == name {
			return &configKey
		}
	}
	lineageKey := parseLineageKey(name)
	if lineageKey != nil {
		return lineageKey
	}
	for _, aliasableCommand := range configdomain.AliasableCommands() {
		key := AliasableCommandToKey(aliasableCommand)
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
	return &Key{
		name: key,
	}
}

// DeprecatedKeys defines the up-to-date counterparts to deprecated configuration settings.
var DeprecatedKeys = map[Key]Key{ //nolint:gochecknoglobals
	KeyDeprecatedCodeHostingDriver:    KeyCodeHostingPlatform,
	KeyDeprecatedMainBranchName:       KeyMainBranch,
	KeyDeprecatedNewBranchPushFlag:    KeyPushNewBranches,
	KeyDeprecatedPerennialBranchNames: KeyPerennialBranches,
	KeyDeprecatedPullBranchStrategy:   KeySyncPerennialStrategy,
	KeyDeprecatedPushVerify:           KeyPushHook,
	KeyDeprecatedSyncStrategy:         KeySyncFeatureStrategy,
}
