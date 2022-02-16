// Package envvars provides helper functions to work with lists of environment variables as provided by `os.Environ`.
package envvars

import (
	"fmt"
	"os"
	"strings"
)

// PrependPath provides a new envvars with the given directory appended to the PATH entry of the given envvars.
// This function assumes there is only one PATH entry.
func PrependPath(envVars []string, directory string) []string {
	for e, envVar := range envVars {
		if strings.HasPrefix(envVar, "PATH=") {
			parts := strings.SplitN(envVar, "=", 2)
			envVars[e] = "PATH=" + directory + string(os.PathListSeparator) + parts[1]
			return envVars
		}
	}
	return append(envVars, fmt.Sprintf("PATH=%s", directory))
}

// Replace provides a new envvars in which the entry with the given key contains the given value instead of its original value.
// If no entry with the given key exists, appends one at the end.
// This function assumes that keys are unique, i.e. no duplicate keys exist.
func Replace(envVars []string, key string, value string) []string {
	prefix := key + "="
	for e, envVar := range envVars {
		if strings.HasPrefix(envVar, prefix) {
			envVars[e] = prefix + value
			return envVars
		}
	}
	return append(envVars, prefix+value)
}
