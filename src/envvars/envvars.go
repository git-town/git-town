// Package envvars provides helper functions to work with lists of environment variables as provided by `os.Environ`.
package envvars

import (
	"fmt"
	"os"
	"strings"
)

// PrependPath provides a new envvars with the given directory appended to the PATH entry of the given envvars.
// This function assumes there is only one PATH entry.
func PrependPath(envvars []string, directory string) []string {
	for i, envVar := range envvars {
		if strings.HasPrefix(envVar, "PATH=") {
			parts := strings.SplitN(envVar, "=", 2)
			envvars[i] = "PATH=" + directory + string(os.PathListSeparator) + parts[1]
			return envvars
		}
	}
	return append(envvars, fmt.Sprintf("PATH=%s", directory))
}

// Replace provides a new envvars in which the entry with the given key contains the given value instead of its original value.
// If no entry with the given key exists, appends one at the end.
// This function assumes that keys are unique, i.e. no duplicate keys exist.
func Replace(envvars []string, key string, value string) []string {
	prefix := key + "="
	for i := range envvars {
		if strings.HasPrefix((envvars)[i], prefix) {
			(envvars)[i] = prefix + value
			return envvars
		}
	}
	return append(envvars, prefix+value)
}
