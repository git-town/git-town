package configdomain

import (
	"encoding/json"
	"fmt"
	"strings"

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
	KeyAliasAppend                    = Key{"alias." + AliasableCommandAppend.String()}       //nolint:gochecknoglobals
	KeyAliasDiffParent                = Key{"alias." + AliasableCommandDiffParent.String()}   //nolint:gochecknoglobals
	KeyAliasHack                      = Key{"alias." + AliasableCommandHack.String()}         //nolint:gochecknoglobals
	KeyAliasKill                      = Key{"alias." + AliasableCommandKill.String()}         //nolint:gochecknoglobals
	KeyAliasPrepend                   = Key{"alias." + AliasableCommandPrepend.String()}      //nolint:gochecknoglobals
	KeyAliasPropose                   = Key{"alias." + AliasableCommandPropose.String()}      //nolint:gochecknoglobals
	KeyAliasRenameBranch              = Key{"alias." + AliasableCommandRenameBranch.String()} //nolint:gochecknoglobals
	KeyAliasRepo                      = Key{"alias." + AliasableCommandRepo.String()}         //nolint:gochecknoglobals
	KeyAliasShip                      = Key{"alias." + AliasableCommandShip.String()}         //nolint:gochecknoglobals
	KeyAliasSync                      = Key{"alias." + AliasableCommandSync.String()}         //nolint:gochecknoglobals
	KeyCodeHostingOriginHostname      = Key{"git-town.code-hosting-origin-hostname"}          //nolint:gochecknoglobals
	KeyCodeHostingPlatform            = Key{"git-town.code-hosting-platform"}                 //nolint:gochecknoglobals
	KeyDeprecatedCodeHostingDriver    = Key{"git-town.code-hosting-driver"}                   //nolint:gochecknoglobals
	KeyDeprecatedMainBranchName       = Key{"git-town.main-branch-name"}                      //nolint:gochecknoglobals
	KeyDeprecatedNewBranchPushFlag    = Key{"git-town.new-branch-push-flag"}                  //nolint:gochecknoglobals
	KeyDeprecatedPerennialBranchNames = Key{"git-town.perennial-branch-names"}                //nolint:gochecknoglobals
	KeyDeprecatedPullBranchStrategy   = Key{"git-town.pull-branch-strategy"}                  //nolint:gochecknoglobals
	KeyDeprecatedPushVerify           = Key{"git-town.push-verify"}                           //nolint:gochecknoglobals
	KeyDeprecatedSyncStrategy         = Key{"git-town.sync-strategy"}                         //nolint:gochecknoglobals
	KeyGiteaToken                     = Key{"git-town.gitea-token"}                           //nolint:gochecknoglobals
	KeyGithubToken                    = Key{"git-town.github-token"}                          //nolint:gochecknoglobals
	KeyGitlabToken                    = Key{"git-town.gitlab-token"}                          //nolint:gochecknoglobals
	KeyMainBranch                     = Key{"git-town.main-branch"}                           //nolint:gochecknoglobals
	KeyOffline                        = Key{"git-town.offline"}                               //nolint:gochecknoglobals
	KeyPerennialBranches              = Key{"git-town.perennial-branches"}                    //nolint:gochecknoglobals
	KeyPushHook                       = Key{"git-town.push-hook"}                             //nolint:gochecknoglobals
	KeyPushNewBranches                = Key{"git-town.push-new-branches"}                     //nolint:gochecknoglobals
	KeyShipDeleteTrackingBranch       = Key{"git-town.ship-delete-remote-branch"}             //nolint:gochecknoglobals
	KeySyncBeforeShip                 = Key{"git-town.sync-before-ship"}                      //nolint:gochecknoglobals
	KeySyncFeatureStrategy            = Key{"git-town.sync-feature-strategy"}                 //nolint:gochecknoglobals
	KeySyncPerennialStrategy          = Key{"git-town.sync-perennial-strategy"}               //nolint:gochecknoglobals
	KeySyncStrategy                   = Key{"git-town.sync-strategy"}                         //nolint:gochecknoglobals
	KeySyncUpstream                   = Key{"git-town.sync-upstream"}                         //nolint:gochecknoglobals
	KeyTestingRemoteURL               = Key{"git-town.testing.remote-url"}                    //nolint:gochecknoglobals
	KeyGitUserEmail                   = Key{"user.email"}                                     //nolint:gochecknoglobals
	KeyGitUserName                    = Key{"user.name"}                                      //nolint:gochecknoglobals
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
	for _, aliasableCommand := range AliasableCommands() {
		if aliasableCommand.Key().String() == name {
			result := aliasableCommand.Key()
			return &result
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
