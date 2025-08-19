package envconfig

import "strings"

// Environment is an immutable representation of all environment variables.
// It allows efficient lookup of environment variables in O(1) time
// by multiple names.
type Environment struct {
	data map[string]string
}

// Get provides the environment variable with the first matching given name.
func (self Environment) Get(name string, alternatives ...string) string {
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

func NewEnvironment(osEnv []string) Environment {
	result := Environment{
		data: map[string]string{},
	}
	for _, entry := range osEnv {
		if name, value, isValid := strings.Cut(entry, "="); isValid {
			result.data[name] = value
		}
	}
	return result
}
