package git

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/git-town/git-town/src/command"
)

// GetVersion indicates whether the needed Git version is installed.
func GetVersion() (int, int, error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	res := command.MustRun("git", "version")
	matches := versionRegexp.FindStringSubmatch(res.OutputSanitized())
	if matches == nil {
		return 0, 0, fmt.Errorf("'git version' returned unexpected output: %q.\nPlease open an issue and supply the output of running 'git version'", res.Output())
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert major version (%v) to int: %w", matches[1], err)
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert minor version (%v) to int: %w", matches[2], err)
	}
	return majorVersion, minorVersion, nil
}
