package envvars

import (
	"os"
	"strings"
)

// PrependPath provides a new envvars with the given directory appended to the PATH entry of the given envvars.
// This function assumes there is only one PATH entry.
func PrependPath(envVars []string, directory string) []string {
	for e, envVar := range envVars {
		if strings.HasPrefix(envVar, "PATH=") {
			_, value, _ := strings.Cut(envVar, "=")
			envVars[e] = "PATH=" + directory + string(os.PathListSeparator) + value
			return envVars
		}
	}
	return append(envVars, "PATH="+directory)
}
