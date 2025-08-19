package envconfig

import "strings"

// EnvVars is an immutable representation of all environment variables.
// It allows efficient lookup of environment variables in O(1) time
// by multiple names.
type EnvVars struct {
	data map[string]string
}

// Get provides the environment variable with the first matching given name.
func (self EnvVars) Get(name string, alternatives ...string) string {
	if result, has := self.data[name]; has {
		return result
	}
	for _, alternative := range alternatives {
		if result, has := self.data[alternative]; has {
			return result
		}
	}
	return ""
}

func NewEnvVars(entries []string) EnvVars {
	result := EnvVars{
		data: map[string]string{},
	}
	for _, entry := range entries {
		if name, value, isValid := strings.Cut(entry, "="); isValid {
			result.data[name] = value
		}
	}
	return result
}
