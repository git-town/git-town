package git

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/util"
)

// EnsureVersionRequirementSatisfied asserts that Git is the needed version or higher
func EnsureVersionRequirementSatisfied() {
	util.Ensure(isVersionRequirementSatisfied(), "Git Town requires Git 2.7.0 or higher")
}

// Helpers

func isVersionRequirementSatisfied() bool {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	res := command.MustRun("git", "version")
	matches := versionRegexp.FindStringSubmatch(res.OutputSanitized())
	if matches == nil {
		log.Fatalf("'git version' returned unexpected output: %q.\nPlease open an issue and supply the output of running 'git version'.", res.Output())
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(fmt.Errorf("cannot convert major version (%v) to int: %w", matches[1], err))
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		panic(fmt.Errorf("cannot convert minor version (%v) to int: %w", matches[2], err))
	}
	return majorVersion == 2 && minorVersion >= 7
}
