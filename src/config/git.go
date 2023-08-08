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
	globalConfigCache map[ConfigKey]string

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfigCache map[ConfigKey]string
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func LoadGit(runner runner, global bool) (map[ConfigKey]string, error) {
	result := map[ConfigKey]string{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := runner.Query("git", cmdArgs...)
	if err != nil {
		return result, err
	}
	if output == "" {
		return result, nil
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		configKey, err := NewConfigKey(key)
		if err != nil {
			return result, err
		}
		result[configKey] = value
	}
	return result, nil
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(runner runner) (Git, error) {
	localConfig, err := LoadGit(runner, false)
	if err != nil {
		return EmptyGit(), err
	}
	globalConfig, err := LoadGit(runner, true)
	if err != nil {
		return EmptyGit(), err
	}
	return Git{
		localConfigCache:  localConfig,
		globalConfigCache: globalConfig,
		runner:            runner,
	}, nil
}

// GlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (g *Git) GlobalConfigValue(key ConfigKey) string {
	return g.globalConfigCache[key]
}

// LocalConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (g *Git) LocalConfigKeysMatching(toMatch string) []ConfigKey {
	result := []ConfigKey{}
	re := regexp.MustCompile(toMatch)
	for key := range g.localConfigCache {
		if re.MatchString(key.String()) {
			result = append(result, key)
		}
	}
	return result
}

// LocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (g *Git) LocalConfigValue(key ConfigKey) string {
	return g.localConfigCache[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (g *Git) LocalOrGlobalConfigValue(key ConfigKey) string {
	local := g.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return g.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (g *Git) Reload() error {
	localConfig, err := LoadGit(g.runner, false)
	if err != nil {
		return err
	}
	globalConfig, err := LoadGit(g.runner, true)
	if err != nil {
		return err
	}
	g.localConfigCache = localConfig
	g.globalConfigCache = globalConfig
	return nil
}

func (g *Git) RemoveGlobalConfigValue(key ConfigKey) (string, error) {
	delete(g.globalConfigCache, key)
	return g.Query("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (g *Git) RemoveLocalConfigValue(key ConfigKey) error {
	delete(g.localConfigCache, key)
	err := g.runner.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (g *Git) SetGlobalConfigValue(key ConfigKey, value string) (string, error) {
	g.globalConfigCache[key] = value
	return g.runner.Query("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (g *Git) SetLocalConfigValue(key ConfigKey, value string) error {
	g.localConfigCache[key] = value
	return g.runner.Run("git", "config", key.String(), value)
}

func EmptyGit() Git {
	return Git{
		runner: emptyRunner{},
	}
}

type runner interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

type emptyRunner struct{}

func (e emptyRunner) Query(string, ...string) (string, error) {
	return "", nil
}

func (e emptyRunner) QueryTrim(string, ...string) (string, error) {
	return "", nil
}

func (e emptyRunner) Run(string, ...string) error {
	return nil
}
