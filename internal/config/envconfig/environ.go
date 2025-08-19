package envconfig

import "strings"

// Environment is an immutable representation of all environment variables.
// It allows lookup by name in O(1) time.
type Environment struct {
	data map[string]string
}

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
