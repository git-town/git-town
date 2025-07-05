package gitconfig

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func LoadSnapshot(runnerquerier subshelldomain.RunnerQuerier, scopeOpt Option[configdomain.ConfigScope], updateOutdated configdomain.UpdateOutdatedSettings) (configdomain.SingleSnapshot, error) {
	snapshot := configdomain.SingleSnapshot{}
	cmdArgs := []string{"config", "-lz"}
	scope, hasScope := scopeOpt.Get()
	if hasScope {
		cmdArgs = append(cmdArgs, scope.GitFlag())
	}
	output, err := runnerquerier.Query("git", cmdArgs...)
	if err != nil || output == "" {
		return snapshot, nil //nolint:nilerr  // Git returns an error if there is no global Git config, assume empty config in this case
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		key, value, _ := strings.Cut(line, "\n")
		configKey, hasConfigKey := configdomain.ParseKey(key).Get()
		if updateOutdated.IsTrue() && hasScope {
			newKey, keyIsDeprecated := configdomain.DeprecatedKeys[configKey]
			if keyIsDeprecated {
				UpdateDeprecatedSetting(runnerquerier, scope, configKey, newKey, value)
				configKey = newKey
			}
			if configKey != configdomain.KeyPerennialBranches && value == "" {
				_ = RemoveConfigValue(runnerquerier, configdomain.ConfigScopeLocal, configKey)
				continue
			}
			if slices.Contains(configdomain.ObsoleteKeys, configKey) {
				_ = RemoveConfigValue(runnerquerier, scope, configKey)
				fmt.Printf(messages.SettingSunsetDeleted, configKey)
				continue
			}
			for _, update := range configdomain.ConfigUpdates {
				if configKey == update.Before.Key && value == update.Before.Value {
					UpdateDeprecatedSetting(runnerquerier, scope, configKey, update.After.Key, update.After.Value)
					configKey = update.After.Key
					value = update.After.Value
				}
			}
			for branchList, branchType := range configdomain.ObsoleteBranchLists {
				if configKey == branchList {
					for _, branch := range strings.Split(value, " ") {
						branchTypeKey := configdomain.Key(configdomain.BranchSpecificKeyPrefix + branch + configdomain.BranchTypeSuffix)
						snapshot[branchTypeKey] = branchType.String()
						_ = SetConfigValue(runnerquerier, configdomain.ConfigScopeLocal, branchTypeKey, branchType.String())
					}
					_ = RemoveConfigValue(runnerquerier, configdomain.ConfigScopeLocal, configKey)
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

func RemoteURL(querier subshelldomain.Querier, remote gitdomain.Remote) Option[string] {
	output, err := querier.Query("git", "remote", "get-url", remote.String())
	if err != nil {
		// NOTE: it's okay to ignore the error here.
		// If we get an error here, we simply don't use the origin remote.
		return None[string]()
	}
	return NewOption(strings.TrimSpace(output))
}

func RemoveConfigValue(runner subshelldomain.Runner, scope configdomain.ConfigScope, key configdomain.Key) error {
	args := []string{"config"}
	if scope == configdomain.ConfigScopeGlobal {
		args = append(args, "--global")
	}
	args = append(args, "--unset", key.String())
	return runner.Run("git", args...)
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func RemoveLocalGitConfiguration(runner subshelldomain.Runner, localSnapshot configdomain.SingleSnapshot) error {
	if err := runner.Run("git", "config", "--remove-section", "git-town"); err != nil {
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
			if err := runner.Run("git", "config", "--unset", key.String()); err != nil {
				return fmt.Errorf(messages.ConfigRemoveError, err)
			}
		}
	}
	return nil
}

// SetConfigValue sets the given configuration setting in the global Git configuration.
func SetConfigValue(runner subshelldomain.Runner, scope configdomain.ConfigScope, key configdomain.Key, value string) error {
	args := []string{"config"}
	if scope == configdomain.ConfigScopeGlobal {
		args = append(args, "--global")
	}
	args = append(args, key.String(), value)
	return runner.Run("git", args...)
}

func UpdateDeprecatedSetting(runner subshelldomain.Runner, scope configdomain.ConfigScope, oldKey, newKey configdomain.Key, value string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingDeprecatedMessage, scope, oldKey, newKey)))
	if err := RemoveConfigValue(runner, scope, oldKey); err != nil {
		fmt.Printf(messages.SettingCannotRemove, scope, oldKey, err)
	}
	if err := SetConfigValue(runner, scope, newKey, value); err != nil {
		fmt.Printf(messages.SettingCannotWrite, scope, newKey, err)
	}
}

// updates a custom Git alias (not set up by Git Town)
func UpdateExternalGitTownAlias(runner subshelldomain.Runner, scope configdomain.ConfigScope, key configdomain.Key, oldValue, newValue string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingDeprecatedValueMessage, scope, key, oldValue, newValue)))
	if err := SetConfigValue(runner, scope, key, newValue); err != nil {
		fmt.Printf(messages.SettingCannotWrite, scope, key, err)
	}
}
