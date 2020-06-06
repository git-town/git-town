package git

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/git-town/git-town/src/command"
)

// CheckVersion indicates whether the needed Git version is installed.
func CheckVersion() error {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	res := command.MustRun("git", "version")
	matches := versionRegexp.FindStringSubmatch(res.OutputSanitized())
	if matches == nil {
		return fmt.Errorf("'git version' returned unexpected output: %q.\nPlease open an issue and supply the output of running 'git version'", res.Output())
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("cannot convert major version (%v) to int: %w", matches[1], err)
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		return fmt.Errorf("cannot convert minor version (%v) to int: %w", matches[2], err)
	}
	if majorVersion*100+minorVersion < 207 {
		// nolint:stylecheck // sentence begins with a proper noun
		return fmt.Errorf("Git Town requires Git 2.7.0 or higher")
	}
	return nil
}
