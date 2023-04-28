package envvars

import "strings"

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
