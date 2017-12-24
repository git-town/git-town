package git

import (
	"strings"

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

// Data returns the map of the data
func (c *ConfigMap) Data() map[string]string {
	c.initialize()
	return c.data
}

// Delete delete the given key
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
	cmdArgs := []string{"git", "config", "-l"}
	if c.global {
		cmdArgs = append(cmdArgs, "--global")
	}
	cmd := command.New(cmdArgs...)
	if cmd.Err() != nil || cmd.Output() == "" {
		return
	}
	for _, line := range strings.Split(cmd.Output(), "\n") {
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]
		c.data[key] = value
	}
	c.initialized = true
}
