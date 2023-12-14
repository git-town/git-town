package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

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
