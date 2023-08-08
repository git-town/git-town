// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

type Key struct {
	name string
}

func (c Key) String() string { return c.name }

var (
	KeyAliasTypeAppend             = Key{"git-town.alias." + AliasTypeAppend.name}
	KeyAliasTypeDiffParent         = Key{"git-town.alias." + AliasTypeDiffParent.name}
	KeyAliasTypeHack               = Key{"git-town.alias." + AliasTypeHack.name}
	KeyAliasTypeKill               = Key{"git-town.alias." + AliasTypeKill.name}
	KeyAliasTypeNewPullRequest     = Key{"git-town.alias." + AliasTypeNewPullRequest.name}
	KeyAliasTypePrepend            = Key{"git-town.alias." + AliasTypePrepend.name}
	KeyAliasTypePruneBranches      = Key{"git-town.alias." + AliasTypePruneBranches.name}
	KeyAliasTypeRenameBranch       = Key{"git-town.alias." + AliasTypeRenameBranch.name}
	KeyAliasTypeRepo               = Key{"git-town.alias." + AliasTypeRepo.name}
	KeyAliasTypeShip               = Key{"git-town.alias." + AliasTypeShip.name}
	KeyAliasTypeSync               = Key{"git-town.alias." + AliasTypeSync.name}
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

func NewAliasKey(aliasType AliasType) Key {
	switch aliasType {
	case AliasTypeAppend:
		return KeyAliasTypeAppend
	case AliasTypeDiffParent:
		return KeyAliasTypeDiffParent
	case AliasTypeHack:
		return KeyAliasTypeHack
	case AliasTypeKill:
		return KeyAliasTypeKill
	case AliasTypeNewPullRequest:
		return KeyAliasTypeNewPullRequest
	case AliasTypePrepend:
		return KeyAliasTypePrepend
	case AliasTypePruneBranches:
		return KeyAliasTypePrepend
	case AliasTypeRenameBranch:
		return KeyAliasTypePruneBranches
	case AliasTypeRepo:
		return KeyAliasTypeRepo
	case AliasTypeShip:
		return KeyAliasTypeShip
	case AliasTypeSync:
		return KeyAliasTypeSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", aliasType))
}
