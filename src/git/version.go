package git

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

// EnsureVersionRequirementSatisfied asserts that Git is the needed version or higher
func EnsureVersionRequirementSatisfied() {
	util.Ensure(isVersionRequirementSatisfied(), "Git Town requires Git 2.7.0 or higher")
}

// Helpers

func isVersionRequirementSatisfied() bool {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	matches := versionRegexp.FindStringSubmatch(command.MustRun("git", "version").OutputSanitized())
	if matches == nil {
		log.Fatal("'git version' returned unexpected output. Please open an issue and supply the output of running 'git version'.")
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
