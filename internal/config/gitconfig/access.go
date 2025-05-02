package gitconfig

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/git-town/git-town/v19/internal/cli/colors"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// Access provides typesafe access to the Git configuration on disk.
type Access struct {
	Runner
}

func (self *Access) Load(scopeOpt Option[configdomain.ConfigScope], updateOutdated bool) (configdomain.SingleSnapshot, error) {
	snapshot := configdomain.SingleSnapshot{}
	cmdArgs := []string{"config", "-lz", "--includes"}
	scope, hasScope := scopeOpt.Get()
	if hasScope {
		cmdArgs = append(cmdArgs, scope.GitFlag())
	}
	output, err := self.Runner.Query("git", cmdArgs...)
	if err != nil {
		// TODO: document why we return nil here when there was an error,
		// or return the error.
		// If returning nil, there is a risk that the initial snapshot is accidentally empty,
		// the final snapshot contains things, and "git town undo" will remove every single Git config value.
		return snapshot, nil //nolint:nilerr
	}
	if output == "" {
		return snapshot, nil
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		configKey, hasConfigKey := configdomain.ParseKey(key).Get()
		if updateOutdated {
			newKey, keyIsDeprecated := configdomain.DeprecatedKeys[configKey]
			if keyIsDeprecated {
				self.UpdateDeprecatedSetting(scope, configKey, newKey, value)
				configKey = newKey
			}
			if configKey != configdomain.KeyPerennialBranches && value == "" {
				_ = self.RemoveLocalConfigValue(configKey)
				continue
			}
			if slices.Contains(configdomain.ObsoleteKeys, configKey) {
				_ = self.RemoveConfigValue(scope, configKey)
				fmt.Printf(messages.SettingSunsetDeleted, configKey)
				continue
			}
			for _, update := range configdomain.ConfigUpdates {
				if configKey == update.Before.Key && value == update.Before.Value {
					self.UpdateDeprecatedSetting(scope, configKey, update.After.Key, update.After.Value)
					configKey = update.After.Key
					value = update.After.Value
				}
			}
			for branchList, branchType := range configdomain.ObsoleteBranchLists {
				if configKey == branchList {
					for _, branch := range strings.Split(value, " ") {
						branchTypeKey := configdomain.Key(configdomain.BranchSpecificKeyPrefix + branch + configdomain.BranchTypeSuffix)
						snapshot[branchTypeKey] = branchType.String()
						_ = self.SetConfigValue(configdomain.ConfigScopeLocal, branchTypeKey, branchType.String())
					}
					_ = self.RemoveLocalConfigValue(configKey)
					fmt.Printf(messages.SettingSunsetBranchList, configKey)
				}
			}
		}
		if hasConfigKey {
			snapshot[configKey] = value
		}
	}
	return snapshot, err
}

func (self *Access) RemoteURL(remote gitdomain.Remote) Option[string] {
	output, err := self.Query("git", "remote", "get-url", remote.String())
	if err != nil {
		// NOTE: it's okay to ignore the error here.
		// If we get an error here, we simply don't use the origin remote.
		return None[string]()
	}
	return NewOption(strings.TrimSpace(output))
}

func (self *Access) RemoveConfigValue(scope configdomain.ConfigScope, key configdomain.Key) error {
	args := []string{"config"}
	if scope == configdomain.ConfigScopeGlobal {
		args = append(args, "--global")
	}
	args = append(args, "--unset", key.String())
	return self.Run("git", args...)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Access) RemoveLocalConfigValue(key configdomain.Key) error {
	return self.Run("git", "config", "--unset", key.String())
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (self *Access) RemoveLocalGitConfiguration(localSnapshot configdomain.SingleSnapshot) error {
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
	for key := range localSnapshot {
		if strings.HasPrefix(key.String(), "git-town-branch.") {
			err = self.Run("git", "config", "--unset", key.String())
			if err != nil {
				return fmt.Errorf(messages.ConfigRemoveError, err)
			}
		}
	}
	return nil
}

// SetConfigValue sets the given configuration setting in the global Git configuration.
func (self *Access) SetConfigValue(scope configdomain.ConfigScope, key configdomain.Key, value string) error {
	args := []string{"config"}
	if scope == configdomain.ConfigScopeGlobal {
		args = append(args, "--global")
	}
	args = append(args, key.String(), value)
	return self.Run("git", args...)
}

func (self *Access) UpdateDeprecatedSetting(scope configdomain.ConfigScope, oldKey, newKey configdomain.Key, value string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingDeprecatedMessage, scope, oldKey, newKey)))
	err := self.RemoveConfigValue(scope, oldKey)
	if err != nil {
		fmt.Printf(messages.SettingCannotRemove, scope, oldKey, err)
	}
	err = self.SetConfigValue(scope, newKey, value)
	if err != nil {
		fmt.Printf(messages.SettingCannotWrite, scope, newKey, err)
	}
}

// updates a custom Git alias (not set up by Git Town)
func (self *Access) UpdateExternalGitTownAlias(scope configdomain.ConfigScope, key configdomain.Key, oldValue, newValue string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingDeprecatedValueMessage, scope, key, oldValue, newValue)))
	err := self.SetConfigValue(scope, key, newValue)
	if err != nil {
		fmt.Printf(messages.SettingCannotWrite, scope, key, err)
	}
}
