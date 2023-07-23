package output

import (
	"strings"
)

// GitCommandsInGitTownOutput provides the Git commands mentioned in the given Git Town output.
func DebugCommandsInGitTownOutput(output string) []string {
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		if lineContainsDebugCommand(line) {
			line = strings.TrimPrefix(line, debugCommandLineBeginning)
			result = append(result, parseDebugLine(line))
		}
	}
	return result
}

// gitCommandLineBeginning contains the first few characters of lines containing Git commands in Git Town output.
const debugCommandLineBeginning = "(debug) " // "\e[1m"

// lineContainsGitTownCommand indicates whether the given line contains a Git Town command.
func lineContainsDebugCommand(line string) bool {
	return strings.HasPrefix(line, debugCommandLineBeginning)
}

// parseLine provides the Git Town command and branch name in the given line.
func parseDebugLine(line string) string {
	return strings.TrimPrefix(line, gitCommandLineBeginning)
}
