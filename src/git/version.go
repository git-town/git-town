package git

import (
	"log"
	"regexp"
	"strconv"

	"github.com/Originate/git-town/src/logs"
	"github.com/Originate/git-town/src/util"
)

// EnsureVersionRequirementSatisfied asserts that Git is the needed version or higher
func EnsureVersionRequirementSatisfied() {
	util.Ensure(isVersionRequirementSatisfied(), "Git Town requires Git 2.7.0 or higher")
}

// Helpers

func isVersionRequirementSatisfied() bool {
	versionRegexp, err := regexp.Compile("git version (\\d+).(\\d+).(\\d+)")
	logs.FatalOnWrap(err, "Error compiling version regular expression")
	matches := versionRegexp.FindStringSubmatch(util.GetCommandOutput("git", "version"))
	if matches == nil {
		log.Fatal("'git version' returned unexpected output. Please open an issue and supply the output of running 'git version'.")
	}
	majorVersion, err := strconv.Atoi(matches[1])
	logs.FatalOnWrap(err, "Error convering major version to int")
	minorVersion, err := strconv.Atoi(matches[2])
	logs.FatalOnWrap(err, "Error convering minor version to int")
	return majorVersion == 2 && minorVersion >= 7
}
