// Package envlist provides helper functions to work with lists of environment variables as provided by `os.Environ`.
package envlist

import (
	"fmt"
	"os"
	"strings"
)

// PrependPath provides a new envList with the given directory appended to the PATH entry of the given envList.
// This function assumes there is only one PATH entry.
func PrependPath(envList []string, directory string) []string {
	for i, envVar := range envList {
		if strings.HasPrefix(envVar, "PATH=") {
			parts := strings.SplitN(envVar, "=", 2)
			envList[i] = "PATH=" + directory + string(os.PathListSeparator) + parts[1]
			return envList
		}
	}
	return append(envList, fmt.Sprintf("PATH=%s", directory))
}

// Replace provides a new envlist in which the entry with the given key contains the given value instead of its original value.
// If no entry with the given key exists, appends one at the end.
// This function assumes that keys are unique, i.e. no duplicate keys exist.
func Replace(envList []string, key string, value string) []string {
	prefix := key + "="
	for i := range envList {
		if strings.HasPrefix((envList)[i], prefix) {
			(envList)[i] = prefix + value
			return envList
		}
	}
	return append(envList, prefix+value)
}
