package gitconfig

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// Access provides typesafe access to the Git configuration on disk.
type Access struct {
	Runner
}

// LoadLocal reads the global Git Town configuration that applies to the entire machine.
func (self *Access) LoadGlobal(updateOutdated bool) (SingleSnapshot, configdomain.PartialConfig, error) {
	return self.load(true, updateOutdated)
}

// LoadLocal reads the Git Town configuration from the local Git's metadata for the current repository.
func (self *Access) LoadLocal(updateOutdated bool) (SingleSnapshot, configdomain.PartialConfig, error) {
	return self.load(false, updateOutdated)
}

func AddKeyToPartialConfig(key Key, value string, config *configdomain.PartialConfig) error {
	if strings.HasPrefix(key.String(), "git-town-branch.") {
		if config.Lineage == nil {
			config.Lineage = configdomain.Lineage{}
		}
		childName := strings.TrimSuffix(strings.TrimPrefix(key.String(), LineageKeyPrefix), LineageKeySuffix)
		child := gitdomain.NewLocalBranchName(childName)
		parent := gitdomain.NewLocalBranchName(value)
		config.Lineage[child] = parent
		return nil
	}
	var err error
	switch key {
	case KeyAliasAppend:
		config.Aliases[configdomain.AliasableCommandAppend] = value
	case KeyAliasCompress:
		config.Aliases[configdomain.AliasableCommandCompress] = value
	case KeyAliasContribute:
		config.Aliases[configdomain.AliasableCommandContribute] = value
	case KeyAliasDiffParent:
		config.Aliases[configdomain.AliasableCommandDiffParent] = value
	case KeyAliasHack:
		config.Aliases[configdomain.AliasableCommandHack] = value
	case KeyAliasKill:
		config.Aliases[configdomain.AliasableCommandKill] = value
	case KeyAliasObserve:
		config.Aliases[configdomain.AliasableCommandObserve] = value
	case KeyAliasPark:
		config.Aliases[configdomain.AliasableCommandPark] = value
	case KeyAliasPrepend:
		config.Aliases[configdomain.AliasableCommandPrepend] = value
	case KeyAliasPropose:
		config.Aliases[configdomain.AliasableCommandPropose] = value
	case KeyAliasRenameBranch:
		config.Aliases[configdomain.AliasableCommandRenameBranch] = value
	case KeyAliasRepo:
		config.Aliases[configdomain.AliasableCommandRepo] = value
	case KeyAliasSetParent:
		config.Aliases[configdomain.AliasableCommandSetParent] = value
	case KeyAliasShip:
		config.Aliases[configdomain.AliasableCommandShip] = value
	case KeyAliasSync:
		config.Aliases[configdomain.AliasableCommandSync] = value
	case KeyContributionBranches:
		config.ContributionBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyHostingOriginHostname:
		config.HostingOriginHostname = configdomain.NewHostingOriginHostnameOption(value)
	case KeyHostingPlatform:
		config.HostingPlatform, err = configdomain.NewHostingPlatformOption(value)
	case KeyGiteaToken:
		config.GiteaToken = configdomain.NewGiteaTokenOption(value)
	case KeyGithubToken:
		config.GitHubToken = configdomain.NewGitHubTokenOption(value)
	case KeyGitlabToken:
		config.GitLabToken = configdomain.NewGitLabTokenOption(value)
	case KeyGitUserEmail:
		config.GitUserEmail = configdomain.NewGitUserEmailOption(value)
	case KeyGitUserName:
		config.GitUserName = configdomain.NewGitUserNameOption(value)
	case KeyMainBranch:
		config.MainBranch = gitdomain.NewLocalBranchNameOption(value)
	case KeyObservedBranches:
		config.ObservedBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyOffline:
		config.Offline, err = configdomain.NewOfflineOption(value, KeyOffline.String())
	case KeyParkedBranches:
		config.ParkedBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyPerennialBranches:
		config.PerennialBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyPerennialRegex:
		config.PerennialRegex = configdomain.NewPerennialRegexOption(value)
	case KeyPushHook:
		var pushHook configdomain.PushHook
		pushHook, err = configdomain.NewPushHook(value, KeyPushHook.String())
		config.PushHook = Some(pushHook)
	case KeyPushNewBranches:
		config.PushNewBranches, err = configdomain.ParsePushNewBranchesOption(value, KeyPushNewBranches.String())
	case KeyShipDeleteTrackingBranch:
		config.ShipDeleteTrackingBranch, err = configdomain.ParseShipDeleteTrackingBranchOption(value, KeyShipDeleteTrackingBranch.String())
	case KeySyncBeforeShip:
		config.SyncBeforeShip, err = configdomain.ParseSyncBeforeShipOption(value, KeySyncBeforeShip.String())
	case KeySyncFeatureStrategy:
		config.SyncFeatureStrategy, err = configdomain.NewSyncFeatureStrategyOption(value)
	case KeySyncPerennialStrategy:
		config.SyncPerennialStrategy, err = configdomain.NewSyncPerennialStrategyOption(value)
	case KeySyncUpstream:
		config.SyncUpstream, err = configdomain.ParseSyncUpstreamOption(value, KeySyncUpstream.String())
	case KeyDeprecatedCodeHostingDriver,
		KeyDeprecatedCodeHostingOriginHostname,
		KeyDeprecatedCodeHostingPlatform,
		KeyDeprecatedMainBranchName,
		KeyDeprecatedNewBranchPushFlag,
		KeyDeprecatedPerennialBranchNames,
		KeyDeprecatedPullBranchStrategy,
		KeyDeprecatedPushVerify,
		KeyDeprecatedShipDeleteRemoteBranch,
		KeyDeprecatedSyncStrategy:
		// deprecated keys were handled before this is reached, they are listed here to check that the switch statement contains all keys
	}
	return err
}

func (self *Access) OriginRemote() string {
	output, err := self.Query("git", "remote", "get-url", gitdomain.RemoteOrigin.String())
	if err != nil {
		// NOTE: it's okay to ignore the error here.
		// If we get an error here, we simply don't use the origin remote.
		return ""
	}
	return strings.TrimSpace(output)
}

func (self *Access) RemoveGlobalConfigValue(key Key) error {
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Access) RemoveLocalConfigValue(key Key) error {
	return self.Run("git", "config", "--unset", key.String())
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (self *Access) RemoveLocalGitConfiguration(lineage configdomain.Lineage) error {
	err := self.Run("git", "config", "--remove-section", "git-town")
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 128 {
				// Git returns exit code 128 when trying to delete a non-existing config section.
				// This is not an error condition in this workflow so we can ignore it here.
				return nil
			}
		}
		return fmt.Errorf(messages.ConfigRemoveError, err)
	}
	for child := range lineage {
		key := fmt.Sprintf("git-town-branch.%s.parent", child)
		err = self.Run("git", "config", "--unset", key)
		if err != nil {
			return fmt.Errorf(messages.ConfigRemoveError, err)
		}
	}
	return nil
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Access) SetGlobalConfigValue(key Key, value string) error {
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Access) SetLocalConfigValue(key Key, value string) error {
	return self.Run("git", "config", key.String(), value)
}

func (self *Access) UpdateDeprecatedGlobalSetting(oldKey, newKey Key, value string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingDeprecatedGlobalMessage, oldKey, newKey)))
	err := self.RemoveGlobalConfigValue(oldKey)
	if err != nil {
		fmt.Printf(messages.SettingGlobalCannotRemove, oldKey, err)
	}
	err = self.SetGlobalConfigValue(newKey, value)
	if err != nil {
		fmt.Printf(messages.SettingGlobalCannotWrite, newKey, err)
	}
}

func (self *Access) UpdateDeprecatedLocalSetting(oldKey, newKey Key, value string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingLocalDeprecatedMessage, oldKey, newKey)))
	err := self.RemoveLocalConfigValue(oldKey)
	if err != nil {
		fmt.Printf(messages.SettingLocalCannotRemove, oldKey, err)
	}
	err = self.SetLocalConfigValue(newKey, value)
	if err != nil {
		fmt.Printf(messages.SettingLocalCannotWrite, newKey, err)
	}
}

func (self *Access) UpdateDeprecatedSetting(oldKey, newKey Key, value string, global bool) {
	if global {
		self.UpdateDeprecatedGlobalSetting(oldKey, newKey, value)
	} else {
		self.UpdateDeprecatedLocalSetting(oldKey, newKey, value)
	}
}

func (self *Access) load(global bool, updateOutdated bool) (SingleSnapshot, configdomain.PartialConfig, error) {
	snapshot := SingleSnapshot{}
	config := configdomain.EmptyPartialConfig()
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := self.Runner.Query("git", cmdArgs...)
	if err != nil {
		return snapshot, config, nil //nolint:nilerr
	}
	if output == "" {
		return snapshot, config, nil
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		configKeyPtr := ParseKey(key)
		if configKeyPtr == nil {
			continue
		}
		configKey := *configKeyPtr
		if updateOutdated {
			newKey, keyIsDeprecated := DeprecatedKeys[configKey]
			if keyIsDeprecated {
				self.UpdateDeprecatedSetting(configKey, newKey, value, global)
				configKey = newKey
			}
			if configKey != KeyPerennialBranches && value == "" {
				_ = self.RemoveLocalConfigValue(configKey)
				continue
			}
		}
		snapshot[configKey] = value
		err := AddKeyToPartialConfig(configKey, value, &config)
		if err != nil {
			return snapshot, config, err
		}
	}
	// verify lineage
	if updateOutdated && config.Lineage != nil {
		for child, parent := range config.Lineage {
			if child == parent {
				fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.ConfigLineageParentIsChild, child)))
				_ = self.RemoveLocalConfigValue(NewParentKey(child))
			}
		}
	}
	return snapshot, config, nil
}
