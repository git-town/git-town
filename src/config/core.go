// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

type ConfigKey struct {
	name string
}

func (c ConfigKey) String() string { return c.name }

var (
	ConfigKeyAliasTypeAppend             = ConfigKey{"git-town.alias." + AliasTypeAppend.name}
	ConfigKeyAliasTypeDiffParent         = ConfigKey{"git-town.alias." + AliasTypeDiffParent.name}
	ConfigKeyAliasTypeHack               = ConfigKey{"git-town.alias." + AliasTypeHack.name}
	ConfigKeyAliasTypeKill               = ConfigKey{"git-town.alias." + AliasTypeKill.name}
	ConfigKeyAliasTypeNewPullRequest     = ConfigKey{"git-town.alias." + AliasTypeNewPullRequest.name}
	ConfigKeyAliasTypePrepend            = ConfigKey{"git-town.alias." + AliasTypePrepend.name}
	ConfigKeyAliasTypePruneBranches      = ConfigKey{"git-town.alias." + AliasTypePruneBranches.name}
	ConfigKeyAliasTypeRenameBranch       = ConfigKey{"git-town.alias." + AliasTypeRenameBranch.name}
	ConfigKeyAliasTypeRepo               = ConfigKey{"git-town.alias." + AliasTypeRepo.name}
	ConfigKeyAliasTypeShip               = ConfigKey{"git-town.alias." + AliasTypeShip.name}
	ConfigKeyAliasTypeSync               = ConfigKey{"git-town.alias." + AliasTypeSync.name}
	ConfigKeyCodeHostingDriver           = ConfigKey{"git-town.code-hosting-driver"}
	ConfigKeyCodeHostingOriginHostname   = ConfigKey{"git-town.code-hosting-origin-hostname"}
	ConfigKeyDeprecatedNewBranchPushFlag = ConfigKey{"git-town.new-branch-push-flag"}
	ConfigKeyDeprecatedPushVerify        = ConfigKey{"git-town.push-verify"}
	ConfigKeyGiteaToken                  = ConfigKey{"git-town.gitea-token"}  //nolint:gosec
	ConfigKeyGithubToken                 = ConfigKey{"git-town.github-token"} //nolint:gosec
	ConfigKeyGitlabToken                 = ConfigKey{"git-town.gitlab-token"} //nolint:gosec
	ConfigKeyMainBranch                  = ConfigKey{"git-town.main-branch-name"}
	ConfigKeyOffline                     = ConfigKey{"git-town.offline"}
	ConfigKeyPerennialBranches           = ConfigKey{"git-town.perennial-branch-names"}
	ConfigKeyPullBranchStrategy          = ConfigKey{"git-town.pull-branch-strategy"}
	ConfigKeyPushHook                    = ConfigKey{"git-town.push-hook"}
	ConfigKeyPushNewBranches             = ConfigKey{"git-town.push-new-branches"}
	ConfigKeyShipDeleteRemoteBranch      = ConfigKey{"git-town.ship-delete-remote-branch"}
	ConfigKeySyncUpstream                = ConfigKey{"git-town.sync-upstream"}
	ConfigKeySyncStrategy                = ConfigKey{"git-town.sync-strategy"}
	ConfigKeyTestingRemoteURL            = ConfigKey{"git-town.testing.remote-url"}
)

var configKeys = []ConfigKey{
	ConfigKeyCodeHostingDriver,
	ConfigKeyCodeHostingOriginHostname,
	ConfigKeyDeprecatedNewBranchPushFlag,
	ConfigKeyDeprecatedPushVerify,
	ConfigKeyGiteaToken,
	ConfigKeyGithubToken,
	ConfigKeyGitlabToken,
	ConfigKeyMainBranch,
	ConfigKeyOffline,
	ConfigKeyPerennialBranches,
	ConfigKeyPullBranchStrategy,
	ConfigKeyPushHook,
	ConfigKeyPushNewBranches,
	ConfigKeyShipDeleteRemoteBranch,
	ConfigKeySyncUpstream,
	ConfigKeySyncStrategy,
	ConfigKeyTestingRemoteURL,
}

func NewConfigKey(value string) (ConfigKey, error) {
	for _, configKey := range configKeys {
		if configKey.name == value {
			return configKey, nil
		}
	}
	return ConfigKeyOffline, fmt.Errorf(messages.ConfigKeyUnknown, value)
}

func NewAliasConfigKey(aliasType AliasType) ConfigKey {
	switch aliasType {
	case AliasTypeAppend:
		return ConfigKeyAliasTypeAppend
	case AliasTypeDiffParent:
		return ConfigKeyAliasTypeDiffParent
	case AliasTypeHack:
		return ConfigKeyAliasTypeHack
	case AliasTypeKill:
		return ConfigKeyAliasTypeKill
	case AliasTypeNewPullRequest:
		return ConfigKeyAliasTypeNewPullRequest
	case AliasTypePrepend:
		return ConfigKeyAliasTypePrepend
	case AliasTypePruneBranches:
		return ConfigKeyAliasTypePrepend
	case AliasTypeRenameBranch:
		return ConfigKeyAliasTypePruneBranches
	case AliasTypeRepo:
		return ConfigKeyAliasTypeRepo
	case AliasTypeShip:
		return ConfigKeyAliasTypeShip
	case AliasTypeSync:
		return ConfigKeyAliasTypeSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", aliasType))
}
