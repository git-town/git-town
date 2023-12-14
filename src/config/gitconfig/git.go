package gitconfig

import (
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// Git provides typesafe access to the Git configuration on disk.
type Git struct {
	Runner
}

func (self *Git) RemoveGlobalConfigValue(key configdomain.Key) error {
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Git) RemoveLocalConfigValue(key configdomain.Key) error {
	err := self.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Git) SetGlobalConfigValue(key configdomain.Key, value string) error {
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Git) SetLocalConfigValue(key configdomain.Key, value string) error {
	return self.Run("git", "config", key.String(), value)
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func (self *Git) LoadCache(global bool) Cache {
	result := Cache{}
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
		if configKey != nil {
			result[*configKey] = value
		}
	}
	return result
}
