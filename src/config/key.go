package config

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/domain"
)

// Key contains all the keys used in Git Town configuration.
type Key struct {
	name string
}

func (c Key) String() string { return c.name }

var (
	KeyAliasAppend                 = Key{"alias." + AliasAppend.name}             //nolint:gochecknoglobals
	KeyAliasDiffParent             = Key{"alias." + AliasDiffParent.name}         //nolint:gochecknoglobals
	KeyAliasHack                   = Key{"alias." + AliasHack.name}               //nolint:gochecknoglobals
	KeyAliasKill                   = Key{"alias." + AliasKill.name}               //nolint:gochecknoglobals
	KeyAliasNewPullRequest         = Key{"alias." + AliasNewPullRequest.name}     //nolint:gochecknoglobals
	KeyAliasPrepend                = Key{"alias." + AliasPrepend.name}            //nolint:gochecknoglobals
	KeyAliasPruneBranches          = Key{"alias." + AliasPruneBranches.name}      //nolint:gochecknoglobals
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

func ParseKey(key string) *Key {
	for _, configKey := range keys {
		if configKey.name == key {
			return &configKey
		}
	}
	lineageKey := ParseLineageKey(key)
	if lineageKey != nil {
		return lineageKey
	}
	return ParseAliasKey(key)
}

func ParseAliasKey(key string) *Key {
	if !strings.HasPrefix(key, "alias.") {
		return nil
	}
	return &Key{
		name: key,
	}
}

func ParseLineageKey(key string) *Key {
	if !strings.HasPrefix(key, "git-town-branch.") || !strings.HasSuffix(key, ".parent") {
		return nil
	}
	return &Key{
		name: key,
	}
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
	case AliasPruneBranches:
		return KeyAliasPruneBranches
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
		name: fmt.Sprintf("git-town-branch.%s.parent", branch),
	}
}
