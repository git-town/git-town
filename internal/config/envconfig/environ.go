package envconfig

import "strings"

// Environment is an immutable representation of all environment variables.
// It allows lookup by name in O(1) time.
type Environment map[string]string

func NewEnvironment(osEnv []string) Environment {
	result := Environment{}
	for _, entry := range osEnv {
		if name, value, isValid := strings.Cut(entry, "="); isValid {
			result[name] = value
		}
	}
	return result
}
