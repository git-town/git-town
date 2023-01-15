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

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func LoadGit(shell run.Shell, global bool) map[string]string {
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
		localConfigCache:  LoadGit(shell, false),
		globalConfigCache: LoadGit(shell, true),
		shell:             shell,
	}
}

// GlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (g *Git) GlobalConfigValue(key string) string {
	return g.globalConfigCache[key]
}

// LocalConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (g *Git) LocalConfigKeysMatching(toMatch string) []string {
	result := []string{}
	re := regexp.MustCompile(toMatch)
	for key := range g.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

// LocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (g *Git) LocalConfigValue(key string) string {
	return g.localConfigCache[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (g *Git) LocalOrGlobalConfigValue(key string) string {
	local := g.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return g.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (g *Git) Reload() {
	g.localConfigCache = LoadGit(g.shell, false)
	g.globalConfigCache = LoadGit(g.shell, true)
}

func (g *Git) RemoveGlobalConfigValue(key string) (*run.Result, error) {
	delete(g.globalConfigCache, key)
	return g.shell.Run("git", "config", "--global", "--unset", key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (g *Git) RemoveLocalConfigValue(key string) error {
	delete(g.localConfigCache, key)
	_, err := g.shell.Run("git", "config", "--unset", key)
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (g *Git) SetGlobalConfigValue(key, value string) (*run.Result, error) {
	g.globalConfigCache[key] = value
	return g.shell.Run("git", "config", "--global", key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (g *Git) SetLocalConfigValue(key, value string) (*run.Result, error) {
	g.localConfigCache[key] = value
	return g.shell.Run("git", "config", key, value)
}
