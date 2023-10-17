package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/domain"
)

// Key contains all the keys used in Git Town configuration.
type Key struct {
	Name string
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (self Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Name)
}

func (self Key) String() string { return self.Name }

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (self *Key) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &self.Name)
}

var (
	KeyAliasAppend                 = Key{"alias." + AliasAppend.name}             //nolint:gochecknoglobals
	KeyAliasDiffParent             = Key{"alias." + AliasDiffParent.name}         //nolint:gochecknoglobals
	KeyAliasHack                   = Key{"alias." + AliasHack.name}               //nolint:gochecknoglobals
	KeyAliasKill                   = Key{"alias." + AliasKill.name}               //nolint:gochecknoglobals
	KeyAliasNewPullRequest         = Key{"alias." + AliasNewPullRequest.name}     //nolint:gochecknoglobals
	KeyAliasPrepend                = Key{"alias." + AliasPrepend.name}            //nolint:gochecknoglobals
	KeyAliasRenameBranch           = Key{"alias." + AliasRenameBranch.name}       //nolint:gochecknoglobals
	KeyAliasRepo                   = Key{"alias." + AliasRepo.name}               //nolint:gochecknoglobals
	KeyAliasShip                   = Key{"alias." + AliasShip.name}               //nolint:gochecknoglobals
	KeyAliasSync                   = Key{"alias." + AliasSync.name}               //nolint:gochecknoglobals
	KeyCodeHostingDriver           = Key{"git-town.code-hosting-driver"}          //nolint:gochecknoglobals
	KeyCodeHostingOriginHostname   = Key{"git-town.code-hosting-origin-hostname"} //nolint:gochecknoglobals
	KeyDeprecatedNewBranchPushFlag = Key{"git-town.new-branch-push-flag"}         //nolint:gochecknoglobals
	KeyDeprecatedPushVerify        = Key{"git-town.push-verify"}                  //nolint:gochecknoglobals
	KeyGiteaToken                  = Key{"git-town.gitea-token"}                  //nolint:gochecknoglobals
	KeyGithubToken                 = Key{"git-town.github-token"}                 //nolint:gochecknoglobals
	KeyGitlabToken                 = Key{"git-town.gitlab-token"}                 //nolint:gochecknoglobals
	KeyMainBranch                  = Key{"git-town.main-branch-name"}             //nolint:gochecknoglobals
	KeyOffline                     = Key{"git-town.offline"}                      //nolint:gochecknoglobals
	KeyPerennialBranches           = Key{"git-town.perennial-branch-names"}       //nolint:gochecknoglobals
	KeyPullBranchStrategy          = Key{"git-town.pull-branch-strategy"}         //nolint:gochecknoglobals
	KeyPushHook                    = Key{"git-town.push-hook"}                    //nolint:gochecknoglobals
	KeyPushNewBranches             = Key{"git-town.push-new-branches"}            //nolint:gochecknoglobals
	KeyShipDeleteRemoteBranch      = Key{"git-town.ship-delete-remote-branch"}    //nolint:gochecknoglobals
	KeySyncUpstream                = Key{"git-town.sync-upstream"}                //nolint:gochecknoglobals
	KeySyncStrategy                = Key{"git-town.sync-strategy"}                //nolint:gochecknoglobals
	KeyTestingRemoteURL            = Key{"git-town.testing.remote-url"}           //nolint:gochecknoglobals
)

var keys = []Key{ //nolint:gochecknoglobals
	KeyCodeHostingDriver,
	KeyCodeHostingOriginHostname,
	KeyDeprecatedNewBranchPushFlag,
	KeyDeprecatedPushVerify,
	KeyGiteaToken,
	KeyGithubToken,
	KeyGitlabToken,
	KeyMainBranch,
	KeyOffline,
	KeyPerennialBranches,
	KeyPullBranchStrategy,
	KeyPushHook,
	KeyPushNewBranches,
	KeyShipDeleteRemoteBranch,
	KeySyncUpstream,
	KeySyncStrategy,
	KeyTestingRemoteURL,
}

func NewAliasKey(aliasType Alias) Key {
	switch aliasType {
	case AliasAppend:
		return KeyAliasAppend
	case AliasDiffParent:
		return KeyAliasDiffParent
	case AliasHack:
		return KeyAliasHack
	case AliasKill:
		return KeyAliasKill
	case AliasNewPullRequest:
		return KeyAliasNewPullRequest
	case AliasPrepend:
		return KeyAliasPrepend
	case AliasRenameBranch:
		return KeyAliasRenameBranch
	case AliasRepo:
		return KeyAliasRepo
	case AliasShip:
		return KeyAliasShip
	case AliasSync:
		return KeyAliasSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", aliasType))
}

func NewParentKey(branch domain.LocalBranchName) Key {
	return Key{
		Name: fmt.Sprintf("git-town-branch.%s.parent", branch),
	}
}

func ParseKey(name string) *Key {
	for _, configKey := range keys {
		if configKey.Name == name {
			return &configKey
		}
	}
	lineageKey := parseLineageKey(name)
	if lineageKey != nil {
		return lineageKey
	}
	return parseAliasKey(name)
}

func parseAliasKey(key string) *Key {
	if !strings.HasPrefix(key, "alias.") {
		return nil
	}
	for _, alias := range Aliases() {
		aliasKey := NewAliasKey(alias)
		if key == aliasKey.Name {
			return &aliasKey
		}
	}
	return nil
}

func parseLineageKey(key string) *Key {
	if !strings.HasPrefix(key, "git-town-branch.") || !strings.HasSuffix(key, ".parent") {
		return nil
	}
	return &Key{
		Name: key,
	}
}
