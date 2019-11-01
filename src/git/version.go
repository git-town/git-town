package git

import (
	"log"
	"regexp"
	"strconv"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
	"github.com/pkg/errors"
)

// EnsureVersionRequirementSatisfied asserts that Git is the needed version or higher
func EnsureVersionRequirementSatisfied() {
	util.Ensure(isVersionRequirementSatisfied(), "Git Town requires Git 2.7.0 or higher")
}

// Helpers

func isVersionRequirementSatisfied() bool {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	matches := versionRegexp.FindStringSubmatch(command.Run("git", "version").OutputSanitized())
	if matches == nil {
		log.Fatal("'git version' returned unexpected output. Please open an issue and supply the output of running 'git version'.")
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(errors.Wrapf(err, "cannot convert major version (%v) to int", matches[1]))
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		panic(errors.Wrapf(err, "cannot convert minor version (%v) to int", matches[2]))
	}
	return majorVersion == 2 && minorVersion >= 7
}
