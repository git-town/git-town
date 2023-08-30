package config

import (
	"regexp"
	"strings"
)

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
	runner

	// globalConfigCache is a cache of the global Git configuration.
	globalConfigCache map[Key]string

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfigCache map[Key]string
}

type runner interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func LoadGit(runner runner, global bool) map[Key]string {
	result := map[Key]string{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := runner.Query("git", cmdArgs...)
	if err != nil {
		return result
	}
	if output == "" {
		return result
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		configKey := ParseKey(key)
		if configKey != nil {
			result[*configKey] = value
		}
	}
	return result
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(runner runner) Git {
	return Git{
		localConfigCache:  LoadGit(runner, false),
		globalConfigCache: LoadGit(runner, true),
		runner:            runner,
	}
}

// GlobalConfig provides the entire global Git configuration
func (g *Git) GlobalConfig() map[Key]string {
	return g.globalConfigCache
}

// GlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (g *Git) GlobalConfigValue(key Key) string {
	return g.globalConfigCache[key]
}

// LocalConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (g *Git) LocalConfigKeysMatching(toMatch string) []Key {
	result := []Key{}
	re := regexp.MustCompile(toMatch)
	for key := range g.localConfigCache {
		if re.MatchString(key.String()) {
			result = append(result, key)
		}
	}
	return result
}

// LocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (g *Git) LocalConfigValue(key Key) string {
	return g.localConfigCache[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (g *Git) LocalOrGlobalConfigValue(key Key) string {
	local := g.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return g.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (g *Git) Reload() {
	g.localConfigCache = LoadGit(g.runner, false)
	g.globalConfigCache = LoadGit(g.runner, true)
}

func (g *Git) RemoveGlobalConfigValue(key Key) (string, error) {
	delete(g.globalConfigCache, key)
	return g.Query("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (g *Git) RemoveLocalConfigValue(key Key) error {
	delete(g.localConfigCache, key)
	err := g.runner.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (g *Git) SetGlobalConfigValue(key Key, value string) (string, error) {
	g.globalConfigCache[key] = value
	return g.runner.Query("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (g *Git) SetLocalConfigValue(key Key, value string) error {
	g.localConfigCache[key] = value
	return g.runner.Run("git", "config", key.String(), value)
}
