package gitconfig

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
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
func (self *Access) LoadGlobal(updateOutdated bool) (configdomain.SingleSnapshot, configdomain.PartialConfig, error) {
	return self.load(true, updateOutdated)
}

// LoadLocal reads the Git Town configuration from the local Git's metadata for the current repository.
func (self *Access) LoadLocal(updateOutdated bool) (configdomain.SingleSnapshot, configdomain.PartialConfig, error) {
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

func (self *Access) RemoveConfigValue(key configdomain.Key, global bool) error {
	if global {
		return self.RemoveGlobalConfigValue(key)
	}
	return self.RemoveLocalConfigValue(key)
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

func (self *Access) load(global bool, updateOutdated bool) (configdomain.SingleSnapshot, configdomain.PartialConfig, error) {
	snapshot := configdomain.SingleSnapshot{}
	cmdArgs := []string{"config", "-lz", "--includes"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := self.Runner.Query("git", cmdArgs...)
	if err != nil {
		return snapshot, configdomain.EmptyPartialConfig(), nil //nolint:nilerr
	}
	if output == "" {
		return snapshot, configdomain.EmptyPartialConfig(), nil
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
			if slices.Contains(configdomain.ObsoleteKeys, configKey) {
				_ = self.RemoveConfigValue(configKey, global)
				fmt.Printf(messages.SettingSunsetDeleted, configKey)
				continue
			}
		}
		snapshot[configKey] = value
	}
	partialConfig, err := configdomain.NewPartialConfigFromSnapshot(snapshot, updateOutdated, self.RemoveLocalConfigValue)
	return snapshot, partialConfig, err
}
