package git

import (
	"log"
	"regexp"
	"strconv"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

// EnsureVersionRequirementSatisfied asserts that Git is the needed version or higher
func EnsureVersionRequirementSatisfied() {
	util.Ensure(isVersionRequirementSatisfied(), "Git Town requires Git 2.7.0 or higher")
}

// Helpers

func isVersionRequirementSatisfied() bool {
	versionRegexp, err := regexp.Compile(`git version (\d+).(\d+).(\d+)`)
	exit.IfWrap(err, "Error compiling version regular expression")
	matches := versionRegexp.FindStringSubmatch(command.New("git", "version").Output())
	if matches == nil {
		log.Fatal("'git version' returned unexpected output. Please open an issue and supply the output of running 'git version'.")
	}
	majorVersion, err := strconv.Atoi(matches[1])
	exit.IfWrap(err, "Error convering major version to int")
	minorVersion, err := strconv.Atoi(matches[2])
	exit.IfWrap(err, "Error convering minor version to int")
	return majorVersion == 2 && minorVersion >= 7
}
