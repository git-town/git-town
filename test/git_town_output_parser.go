package test

import (
	"strings"
)

// GitTownCommandsInOutput provides the Git Town commands in the given output string
func GitTownCommandsInOutput(output string) []string {
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if lineContainsGitTownCommand(line) {
			command, _ := parseLine(line)
			result = append(result, command)
		}
	}
	return result
}

var linePrefix = "\x1b[1m" // "\e[1m"

// lineContainsGitTownCommand indicates whether the given line contains a Git Town command
func lineContainsGitTownCommand(line string) bool {
	return strings.HasPrefix(line, linePrefix)
}

// parseLine returns the Git Town command and branchname in the given line
func parseLine(line string) (command, branch string) {
	// NOTE: implementing this without regex
	// because the regex has gotten very complex and hard to maintain

	// remove the bold formatting
	line = strings.Replace(line, linePrefix, "", 1)

	// extract branch name if it exists
	branchName := ""
	if line[0] == '[' {
		// line contains a branch name
		line = line[1:len(line)] // remove the leading "["
		parts := strings.SplitN(line, "]", 2)
		branchName = parts[0]
		line = parts[1]
	}

	return strings.TrimSpace(line), branchName
}
