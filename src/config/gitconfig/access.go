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

// Note: this exists here and not as a method of PartialConfig to avoid circular dependencies
func (self *Access) AddValueToPartialConfig(key configdomain.Key, value string, config *configdomain.PartialConfig) error {
	if strings.HasPrefix(key.String(), configdomain.LineageKeyPrefix) {
		childName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(key.String(), configdomain.LineageKeyPrefix), configdomain.LineageKeySuffix))
		if childName == "" {
			// empty lineage entries are invalid --> delete it
			return self.RemoveLocalConfigValue(key)
		}
		child := gitdomain.NewLocalBranchName(childName)
		value = strings.TrimSpace(value)
		if value == "" {
			// empty lineage entries are invalid --> delete it
			return self.RemoveLocalConfigValue(key)
		}
		parent := gitdomain.NewLocalBranchName(value)
		config.Lineage.Add(child, parent)
		return nil
	}
	var err error
	switch key {
	case configdomain.KeyAliasAppend:
		config.Aliases[configdomain.AliasableCommandAppend] = value
	case configdomain.KeyAliasCompress:
		config.Aliases[configdomain.AliasableCommandCompress] = value
	case configdomain.KeyAliasContribute:
		config.Aliases[configdomain.AliasableCommandContribute] = value
	case configdomain.KeyAliasDiffParent:
		config.Aliases[configdomain.AliasableCommandDiffParent] = value
	case configdomain.KeyAliasHack:
		config.Aliases[configdomain.AliasableCommandHack] = value
	case configdomain.KeyAliasKill:
		config.Aliases[configdomain.AliasableCommandKill] = value
	case configdomain.KeyAliasObserve:
		config.Aliases[configdomain.AliasableCommandObserve] = value
	case configdomain.KeyAliasPark:
		config.Aliases[configdomain.AliasableCommandPark] = value
	case configdomain.KeyAliasPrepend:
		config.Aliases[configdomain.AliasableCommandPrepend] = value
	case configdomain.KeyAliasPropose:
		config.Aliases[configdomain.AliasableCommandPropose] = value
	case configdomain.KeyAliasRenameBranch:
		config.Aliases[configdomain.AliasableCommandRenameBranch] = value
	case configdomain.KeyAliasRepo:
		config.Aliases[configdomain.AliasableCommandRepo] = value
	case configdomain.KeyAliasSetParent:
		config.Aliases[configdomain.AliasableCommandSetParent] = value
	case configdomain.KeyAliasShip:
		config.Aliases[configdomain.AliasableCommandShip] = value
	case configdomain.KeyAliasSync:
		config.Aliases[configdomain.AliasableCommandSync] = value
	case configdomain.KeyContributionBranches:
		config.ContributionBranches = gitdomain.ParseLocalBranchNames(value)
	case configdomain.KeyCreatePrototypeBranches:
		var createPrototypeBranches configdomain.CreatePrototypeBranches
		createPrototypeBranches, err = configdomain.NewCreatePrototypeBranches(value, configdomain.KeyPrototypeBranches.String())
		config.CreatePrototypeBranches = Some(createPrototypeBranches)
	case configdomain.KeyHostingOriginHostname:
		config.HostingOriginHostname = configdomain.NewHostingOriginHostnameOption(value)
	case configdomain.KeyHostingPlatform:
		config.HostingPlatform, err = configdomain.NewHostingPlatformOption(value)
	case configdomain.KeyGiteaToken:
		config.GiteaToken = configdomain.NewGiteaTokenOption(value)
	case configdomain.KeyGithubToken:
		config.GitHubToken = configdomain.NewGitHubTokenOption(value)
	case configdomain.KeyGitlabToken:
		config.GitLabToken = configdomain.NewGitLabTokenOption(value)
	case configdomain.KeyGitUserEmail:
		config.GitUserEmail = configdomain.NewGitUserEmailOption(value)
	case configdomain.KeyGitUserName:
		config.GitUserName = configdomain.NewGitUserNameOption(value)
	case configdomain.KeyMainBranch:
		config.MainBranch = gitdomain.NewLocalBranchNameOption(value)
	case configdomain.KeyObservedBranches:
		config.ObservedBranches = gitdomain.ParseLocalBranchNames(value)
	case configdomain.KeyOffline:
		config.Offline, err = configdomain.NewOfflineOption(value, configdomain.KeyOffline.String())
	case configdomain.KeyParkedBranches:
		config.ParkedBranches = gitdomain.ParseLocalBranchNames(value)
	case configdomain.KeyPerennialBranches:
		config.PerennialBranches = gitdomain.ParseLocalBranchNames(value)
	case configdomain.KeyPerennialRegex:
		config.PerennialRegex = configdomain.NewPerennialRegexOption(value)
	case configdomain.KeyPrototypeBranches:
		config.PrototypeBranches = gitdomain.ParseLocalBranchNames(value)
	case configdomain.KeyPushHook:
		var pushHook configdomain.PushHook
		pushHook, err = configdomain.NewPushHook(value, configdomain.KeyPushHook.String())
		config.PushHook = Some(pushHook)
	case configdomain.KeyPushNewBranches:
		config.PushNewBranches, err = configdomain.ParsePushNewBranchesOption(value, configdomain.KeyPushNewBranches.String())
	case configdomain.KeyShipDeleteTrackingBranch:
		config.ShipDeleteTrackingBranch, err = configdomain.ParseShipDeleteTrackingBranchOption(value, configdomain.KeyShipDeleteTrackingBranch.String())
	case configdomain.KeySyncBeforeShip:
		config.SyncBeforeShip, err = configdomain.ParseSyncBeforeShipOption(value, configdomain.KeySyncBeforeShip.String())
	case configdomain.KeySyncFeatureStrategy:
		config.SyncFeatureStrategy, err = configdomain.NewSyncFeatureStrategyOption(value)
	case configdomain.KeySyncPerennialStrategy:
		config.SyncPerennialStrategy, err = configdomain.NewSyncPerennialStrategyOption(value)
	case configdomain.KeySyncUpstream:
		config.SyncUpstream, err = configdomain.ParseSyncUpstreamOption(value, configdomain.KeySyncUpstream.String())
	case configdomain.KeyDeprecatedCodeHostingDriver,
		configdomain.KeyDeprecatedCodeHostingOriginHostname,
		configdomain.KeyDeprecatedCodeHostingPlatform,
		configdomain.KeyDeprecatedMainBranchName,
		configdomain.KeyDeprecatedNewBranchPushFlag,
		configdomain.KeyDeprecatedPerennialBranchNames,
		configdomain.KeyDeprecatedPullBranchStrategy,
		configdomain.KeyDeprecatedPushVerify,
		configdomain.KeyDeprecatedShipDeleteRemoteBranch,
		configdomain.KeyDeprecatedSyncStrategy:
		// deprecated keys were handled before this is reached, they are listed here to check that the switch statement contains all keys
	}
	return err
}

// LoadLocal reads the global Git Town configuration that applies to the entire machine.
func (self *Access) LoadGlobal(updateOutdated bool) (SingleSnapshot, configdomain.PartialConfig, error) {
	return self.load(true, updateOutdated)
}

// LoadLocal reads the Git Town configuration from the local Git's metadata for the current repository.
func (self *Access) LoadLocal(updateOutdated bool) (SingleSnapshot, configdomain.PartialConfig, error) {
	return self.load(false, updateOutdated)
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

func (self *Access) RemoveGlobalConfigValue(key configdomain.Key) error {
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Access) RemoveLocalConfigValue(key configdomain.Key) error {
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
	for _, entry := range lineage.Entries() {
		key := fmt.Sprintf("git-town-branch.%s.parent", entry.Child)
		err = self.Run("git", "config", "--unset", key)
		if err != nil {
			return fmt.Errorf(messages.ConfigRemoveError, err)
		}
	}
	return nil
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Access) SetGlobalConfigValue(key configdomain.Key, value string) error {
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Access) SetLocalConfigValue(key configdomain.Key, value string) error {
	return self.Run("git", "config", key.String(), value)
}

func (self *Access) UpdateDeprecatedGlobalSetting(oldKey, newKey configdomain.Key, value string) {
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

func (self *Access) UpdateDeprecatedLocalSetting(oldKey, newKey configdomain.Key, value string) {
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

func (self *Access) UpdateDeprecatedSetting(oldKey, newKey configdomain.Key, value string, global bool) {
	if global {
		self.UpdateDeprecatedGlobalSetting(oldKey, newKey, value)
	} else {
		self.UpdateDeprecatedLocalSetting(oldKey, newKey, value)
	}
}

func (self *Access) load(global bool, updateOutdated bool) (SingleSnapshot, configdomain.PartialConfig, error) {
	snapshot := SingleSnapshot{}
	config := configdomain.EmptyPartialConfig()
	cmdArgs := []string{"config", "-lz", "--includes"}
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
		configKey, hasConfigKey := configdomain.ParseKey(key).Get()
		if !hasConfigKey {
			continue
		}
		if updateOutdated {
			newKey, keyIsDeprecated := configdomain.DeprecatedKeys[configKey]
			if keyIsDeprecated {
				self.UpdateDeprecatedSetting(configKey, newKey, value, global)
				configKey = newKey
			}
			if configKey != configdomain.KeyPerennialBranches && value == "" {
				_ = self.RemoveLocalConfigValue(configKey)
				continue
			}
		}
		snapshot[configKey] = value
		err := self.AddValueToPartialConfig(configKey, value, &config)
		if err != nil {
			return snapshot, config, err
		}
	}
	// verify lineage
	if updateOutdated {
		for _, entry := range config.Lineage.Entries() {
			if entry.Child == entry.Parent {
				fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.ConfigLineageParentIsChild, entry.Child)))
				_ = self.RemoveLocalConfigValue(configdomain.NewParentKey(entry.Child))
			}
		}
	}
	return snapshot, config, nil
}
