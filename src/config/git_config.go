package config

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/run"
)

// Config manages the Git Town configuration
// stored in Git metadata in the given local repo and the global Git configuration.
// This class manages which config values are stored in local vs global settings.
type gitConfig struct {
	// globalConfigCache is a cache of the global Git configuration.
	globalConfigCache map[string]string

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfigCache map[string]string

	// for running shell commands
	shell run.Shell
}

// loadGitConfig provides the Git configuration from the given directory or the global one if the global flag is set.
func loadGitConfig(shell run.Shell, global bool) map[string]string {
	result := map[string]string{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	res, err := shell.Run("git", cmdArgs...)
	if err != nil {
		return result
	}
	output := res.Output()
	if output == "" {
		return result
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		result[key] = value
	}
	return result
}

// globalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *gitConfig) globalConfigValue(key string) string {
	return c.globalConfigCache[key]
}

// localConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *gitConfig) localConfigValue(key string) string {
	return c.localConfigCache[key]
}

// Reload refreshes the cached configuration information.
func (c *gitConfig) Reload() {
	c.localConfigCache = loadGitConfig(c.shell, false)
	c.globalConfigCache = loadGitConfig(c.shell, true)
}

func (c *gitConfig) removeGlobalConfigValue(key string) (*run.Result, error) {
	delete(c.globalConfigCache, key)
	return c.shell.Run("git", "config", "--global", "--unset", key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (c *gitConfig) removeLocalConfigValue(key string) error {
	delete(c.localConfigCache, key)
	_, err := c.shell.Run("git", "config", "--unset", key)
	return err
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (c *gitConfig) RemoveLocalGitConfiguration() error {
	result, err := c.shell.Run("git", "config", "--remove-section", "git-town")
	if err != nil {
		if result.ExitCode() == 128 {
			// Git returns exit code 128 when trying to delete a non-existing config section.
			// This is not an error condition in this workflow so we can ignore it here.
			return nil
		}
		return fmt.Errorf("unexpected error while removing the 'git-town' section from the Git configuration: %w", err)
	}
	return nil
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (c *gitConfig) SetGlobalConfigValue(key, value string) (*run.Result, error) {
	c.globalConfigCache[key] = value
	return c.shell.Run("git", "config", "--global", key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (c *gitConfig) SetLocalConfigValue(key, value string) (*run.Result, error) {
	c.localConfigCache[key] = value
	return c.shell.Run("git", "config", key, value)
}
