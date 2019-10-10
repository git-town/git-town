package git

import (
	"regexp"
)

// ConfigCache represents the data from a call to
// `git config -l` or `git config -l --global`
// TODO: move the initialize method to the place where this is used
type ConfigCache struct {
	data map[string]string
}

// NewConfigCache returns a new config map
func NewConfigCache(global bool) ConfigCache {
	return ConfigCache{
		data: map[string]string{},
	}
}

// KeysMatching returns the keys that match the given regexp
func (c *ConfigCache) KeysMatching(re *regexp.Regexp) (result []string) {
	for key := range c.data {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return
}

// Delete deletes the given key
func (c *ConfigCache) Delete(key string) {
	delete(c.data, key)
}

// Get returns the value for the given key
func (c *ConfigCache) Get(key string) string {
	return c.data[key]
}

// Set updates a key/value pair of the data
func (c *ConfigCache) Set(key, value string) {
	c.data[key] = value
}
