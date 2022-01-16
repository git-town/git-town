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
	for i := range envList {
		if strings.HasPrefix(envList[i], "PATH=") {
			parts := strings.SplitN(envList[i], "=", 2)
			parts[1] = directory + string(os.PathListSeparator) + parts[1]
			envList[i] = strings.Join(parts, "=")
			return envList
		}
	}
	return append(envList, fmt.Sprintf("PATH=%s", directory))
}

// Replace provides a new envlist in which the entry with the given key contains the given value instead of its original value.
// If no entry with the given key exists, appends one at the end.
// This function assumes that keys are unique, i.e. no duplicate keys exist.
func Replace(envList []string, key string, value string) []string {
	for i := range envList {
		if strings.HasPrefix((envList)[i], key+"=") {
			(envList)[i] = fmt.Sprintf("%s=%s", key, value)
			return envList
		}
	}
	return append(envList, fmt.Sprintf("%s=%s", key, value))
}
