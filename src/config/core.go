// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

// Key contains all the keys used in Git Town configuration.
type Key struct {
	name string
}

func (c Key) String() string { return c.name }

var (
	KeyAliasAppend                 = Key{"git-town.alias." + AliasAppend.name}
	KeyAliasTypeParent             = Key{"git-town.alias." + AliasDiffParent.name}
	KeyAliasHack                   = Key{"git-town.alias." + AliasHack.name}
	KeyAliasKill                   = Key{"git-town.alias." + AliasKill.name}
	KeyAliasNewPullRequest         = Key{"git-town.alias." + AliasNewPullRequest.name}
	KeyAliasPrepend                = Key{"git-town.alias." + AliasPrepend.name}
	KeyAliasPruneBranches          = Key{"git-town.alias." + AliasPruneBranches.name}
	KeyAliasRenameBranch           = Key{"git-town.alias." + AliasRenameBranch.name}
	KeyAliasRepo                   = Key{"git-town.alias." + AliasRepo.name}
	KeyAliasShip                   = Key{"git-town.alias." + AliasShip.name}
	KeyAliasSync                   = Key{"git-town.alias." + AliasSync.name}
	KeyCodeHostingDriver           = Key{"git-town.code-hosting-driver"}
	KeyCodeHostingOriginHostname   = Key{"git-town.code-hosting-origin-hostname"}
	KeyDeprecatedNewBranchPushFlag = Key{"git-town.new-branch-push-flag"}
	KeyDeprecatedPushVerify        = Key{"git-town.push-verify"}
	KeyGiteaToken                  = Key{"git-town.gitea-token"}  //nolint:gosec
	KeyGithubToken                 = Key{"git-town.github-token"} //nolint:gosec
	KeyGitlabToken                 = Key{"git-town.gitlab-token"} //nolint:gosec
	KeyMainBranch                  = Key{"git-town.main-branch-name"}
	KeyOffline                     = Key{"git-town.offline"}
	KeyPerennialBranches           = Key{"git-town.perennial-branch-names"}
	KeyPullBranchStrategy          = Key{"git-town.pull-branch-strategy"}
	KeyPushHook                    = Key{"git-town.push-hook"}
	KeyPushNewBranches             = Key{"git-town.push-new-branches"}
	KeyShipDeleteRemoteBranch      = Key{"git-town.ship-delete-remote-branch"}
	KeySyncUpstream                = Key{"git-town.sync-upstream"}
	KeySyncStrategy                = Key{"git-town.sync-strategy"}
	KeyTestingRemoteURL            = Key{"git-town.testing.remote-url"}
)

var keys = []Key{
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

func NewKey(value string) (Key, error) {
	for _, configKey := range keys {
		if configKey.name == value {
			return configKey, nil
		}
	}
	return KeyOffline, fmt.Errorf(messages.ConfigKeyUnknown, value)
}

func NewAliasKey(aliasType Alias) Key {
	switch aliasType {
	case AliasAppend:
		return KeyAliasAppend
	case AliasDiffParent:
		return KeyAliasTypeParent
	case AliasHack:
		return KeyAliasHack
	case AliasKill:
		return KeyAliasKill
	case AliasNewPullRequest:
		return KeyAliasNewPullRequest
	case AliasPrepend:
		return KeyAliasPrepend
	case AliasPruneBranches:
		return KeyAliasPrepend
	case AliasRenameBranch:
		return KeyAliasPruneBranches
	case AliasRepo:
		return KeyAliasRepo
	case AliasShip:
		return KeyAliasShip
	case AliasSync:
		return KeyAliasSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", aliasType))
}
