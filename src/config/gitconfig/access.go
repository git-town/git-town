package gitconfig

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// Access provides typesafe access to the Git configuration on disk.
type Access struct {
	Runner
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func (self *Access) LoadCache(global bool) SingleCache {
	result := SingleCache{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := self.Runner.Query("git", cmdArgs...)
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
		configKey := configdomain.ParseKey(key)
		if configKey == nil {
			continue
		}
		newKey, keyIsDeprecated := configdomain.DeprecatedKeys[*configKey]
		if keyIsDeprecated {
			self.UpdateDeprecatedSetting(*configKey, newKey, value, global)
			configKey = &newKey
		}
		result[*configKey] = value
	}
	return result
}

func (self *Access) RemoveGlobalConfigValue(key configdomain.Key) error {
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Access) RemoveLocalConfigValue(key configdomain.Key) error {
	err := self.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Access) SetGlobalConfigValue(key configdomain.Key, value string) error {
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Access) SetLocalConfigValue(key configdomain.Key, value string) error {
	return self.Run("git", "config", key.String(), value)
}

func (self *Access) UpdateDeprecatedSetting(oldKey, newKey configdomain.Key, value string, global bool) {
	fmt.Printf("I found the deprecated local setting %q.\n", oldKey)
	fmt.Printf("I am upgrading this setting to the new format %q.\n", newKey)
	if global {
		err := self.RemoveGlobalConfigValue(oldKey)
		if err != nil {
			fmt.Printf("ERROR: cannot remove global Git setting %q: %v", oldKey, err)
		}
		err = self.SetGlobalConfigValue(newKey, value)
		if err != nil {
			fmt.Printf("ERROR: cannot write global Git setting %q: %v", newKey, err)
		}
	} else {
		err := self.RemoveLocalConfigValue(oldKey)
		if err != nil {
			fmt.Printf("ERROR: cannot remove local Git setting %q: %v", oldKey, err)
		}
		err = self.SetLocalConfigValue(newKey, value)
		if err != nil {
			fmt.Printf("ERROR: cannot write local Git setting %q: %v", newKey, err)
		}
	}
}
