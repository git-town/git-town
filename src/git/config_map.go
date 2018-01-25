package git

import (
	"regexp"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/command"
)

// ConfigMap represents the data from a call to
// `git config -l` or `git config -l --global`
type ConfigMap struct {
	data        map[string]string
	global      bool
	initialized bool
}

// NewConfigMap returns a new config map
func NewConfigMap(global bool) *ConfigMap {
	return &ConfigMap{
		data:        map[string]string{},
		global:      global,
		initialized: false,
	}
}

// KeysMatching returns the keys that match the given regexp
func (c *ConfigMap) KeysMatching(re *regexp.Regexp) (result []string) {
	c.initialize()
	for key := range c.data {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return
}

// Delete deletes the given key
func (c *ConfigMap) Delete(key string) {
	c.initialize()
	delete(c.data, key)
}

// Get returns the value for the given key
func (c *ConfigMap) Get(key string) string {
	c.initialize()
	return c.data[key]
}

// Set updates a key/value pair of the data
func (c *ConfigMap) Set(key, value string) {
	c.initialize()
	c.data[key] = value
}

// Reset resets the configuration map
func (c *ConfigMap) Reset() {
	c.initialized = false
}

// Helpers

func (c *ConfigMap) initialize() {
	if c.initialized {
		return
	}
	cmdArgs := []string{"git", "config", "-lz"}
	if c.global {
		cmdArgs = append(cmdArgs, "--global")
	}
	cmd := command.New(cmdArgs...)
	if cmd.Err() != nil && strings.Contains(cmd.Output(), "No such file or directory") {
		return
	}
	exit.If(cmd.Err())
	if cmd.Output() == "" {
		return
	}
	for _, line := range strings.Split(cmd.Output(), "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		c.data[key] = value
	}
	c.initialized = true
}
