package output

import (
	"strings"
)

// GitCommandsInGitTownOutput provides the Git commands mentioned in the given Git Town output.
func DebugCommandsInGitTownOutput(output string) []string {
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
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
	// NOTE: implementing this without regex because the regex has gotten very complex and hard to maintain
	// remove the color codes at the beginning
	line = strings.Replace(line, gitCommandLineBeginning, "", 1)
	// extract branch name if it exists
	if line[0] == '[' {
		// line contains a branch name
		line = line[1:] // remove the leading "["
		parts := strings.SplitN(line, "]", 2)
		line = parts[1]
	}
	return strings.TrimSpace(line)
}
