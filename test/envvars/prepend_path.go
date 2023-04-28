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
