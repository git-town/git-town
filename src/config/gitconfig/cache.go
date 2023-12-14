package gitconfig

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"golang.org/x/exp/maps"
)

type Cache map[configdomain.Key]string

// Clone provides a copy of this GitConfiguration instance.
func (self Cache) Clone() Cache {
	result := Cache{}
	maps.Copy(result, self)
	return result
}

// KeysMatching provides the keys in this GitConfigCache that match the given regex.
func (self Cache) KeysMatching(pattern string) []configdomain.Key {
	result := []configdomain.Key{}
	re := regexp.MustCompile(pattern)
	for key := range self {
		if re.MatchString(key.String()) {
			result = append(result, key)
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].String() < result[b].String() })
	return result
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func LoadGitConfigCache(runner Runner, global bool) (config configdomain.PartialConfig, cache Cache) {
	cache = Cache{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := runner.Query("git", cmdArgs...)
	if err != nil {
		return config, cache
	}
	if output == "" {
		return config, cache
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		configKey := configdomain.ParseKey(key)
		if configKey == nil {
			// this part of the Git configuration is not a Git Town specific config setting --> ignore
			continue
		}
		newConfigKey := updateDeprecatedKey(*configKey, value, runner)
		if newConfigKey != nil {
			configKey = newConfigKey
		}
		if config.Add(*configKey, value) {
			continue
		}
		cache[*configKey] = value
	}
	return config, cache
}

// updateDeprecatedKey updates the given deprecated config key to its up-to-date counterpart and returns the latter.
func updateDeprecatedKey(key configdomain.Key, value string, global bool, runner Runner) configdomain.Key {
	newKey, isDeprecated := configdomain.DeprecatedKeys[key]
	if !isDeprecated {
		return key
	}
	if global {
		updateDeprecatedGlobalSetting(key, newKey, value, runner)
	}
	return newKey
}

func updateDeprecatedGlobalSetting(deprecatedKey, newKey configdomain.Key, value string, runner Runner) error {
	fmt.Printf("I found the deprecated global setting %q.\n", deprecatedKey)
	fmt.Printf("I am upgrading this setting to the new format %q.\n", newKey)
	err := RemoveGlobalConfigValue(deprecatedKey)
	if err != nil {
		return err
	}
	err = self.SetGlobalConfigValue(newKey, deprecatedSetting)
	return err
	return nil
}

func (self *GitTown) updateDeprecatedLocalSetting(deprecatedKey, newKey configdomain.Key) error {
	deprecatedSetting := self.LocalConfigValue(deprecatedKey)
	if deprecatedSetting != "" {
		fmt.Printf("I found the deprecated local setting %q.\n", deprecatedKey)
		fmt.Printf("I am upgrading this setting to the new format %q.\n", newKey)
		err := self.RemoveLocalConfigValue(deprecatedKey)
		if err != nil {
			return err
		}
		err = self.SetLocalConfigValue(newKey, deprecatedSetting)
		return err
	}
	return nil
}

func (self *GitTown) updateDeprecatedSetting(deprecatedKey, newKey configdomain.Key) error {
	err := self.updateDeprecatedLocalSetting(deprecatedKey, newKey)
	if err != nil {
		return err
	}
	return self.updateDeprecatedGlobalSetting(deprecatedKey, newKey)
}
