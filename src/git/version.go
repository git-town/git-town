package git

import (
	"log"
	"regexp"
	"strconv"

	"github.com/Originate/git-town/src/util"
)

// EnsureVersionRequirementSatisfied asserts that Git is the needed version or higher
func EnsureVersionRequirementSatisfied() {
	util.Ensure(isVersionRequirementSatisfied(), "Git Town requires Git 2.7.0 or higher")
}

// Helpers

func isVersionRequirementSatisfied() bool {
	versionRegexp, err := regexp.Compile("git version (\\d+).(\\d+).(\\d+)")
	if err != nil {
		log.Fatal("Error compiling version regular expression: ", err)
	}
	matches := versionRegexp.FindStringSubmatch(util.GetCommandOutput("git", "version"))
	if matches == nil {
		log.Fatal("'git version' returned unexpected output. Please open an issue and supply the output of running 'git version'.")
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal("Error convering major version to int:", err)
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Fatal("Error convering minor version to int:", err)
	}
	return majorVersion == 2 && minorVersion >= 7
}
