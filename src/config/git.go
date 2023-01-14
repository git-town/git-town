package config

import (
	"regexp"
	"strings"

	"github.com/git-town/git-town/v7/src/run"
)

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
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

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(shell run.Shell) Git {
	return Git{
		localConfigCache:  loadGitConfig(shell, false),
		globalConfigCache: loadGitConfig(shell, true),
		shell:             shell,
	}
}

// GlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (gc *Git) GlobalConfigValue(key string) string {
	return gc.globalConfigCache[key]
}

// LocalConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (gc *Git) LocalConfigKeysMatching(toMatch string) []string {
	result := []string{}
	re := regexp.MustCompile(toMatch)
	for key := range gc.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

// LocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (gc *Git) LocalConfigValue(key string) string {
	return gc.localConfigCache[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (gc *Git) LocalOrGlobalConfigValue(key string) string {
	local := gc.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return gc.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (gc *Git) Reload() {
	gc.localConfigCache = loadGitConfig(gc.shell, false)
	gc.globalConfigCache = loadGitConfig(gc.shell, true)
}

func (gc *Git) removeGlobalConfigValue(key string) (*run.Result, error) {
	delete(gc.globalConfigCache, key)
	return gc.shell.Run("git", "config", "--global", "--unset", key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (gc *Git) removeLocalConfigValue(key string) error {
	delete(gc.localConfigCache, key)
	_, err := gc.shell.Run("git", "config", "--unset", key)
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (gc *Git) SetGlobalConfigValue(key, value string) (*run.Result, error) {
	gc.globalConfigCache[key] = value
	return gc.shell.Run("git", "config", "--global", key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (gc *Git) SetLocalConfigValue(key, value string) (*run.Result, error) {
	gc.localConfigCache[key] = value
	return gc.shell.Run("git", "config", key, value)
}
