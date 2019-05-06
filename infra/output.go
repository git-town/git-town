package infra

import (
	"strings"
)

// CommandsInOutput returns the commands in the given output string
func CommandsInOutput(output string) []string {
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			if lineContainsCommand(line) {
				command, _ := parseLine(line)
				result = append(result, command)
			}
		}
	}
	return result
}

var linePrefix = "\x1b[1m" // "\e[1m"

// lineContainsCommand returns whether the given line contains the given command
func lineContainsCommand(line string) bool {
	return strings.HasPrefix(line, linePrefix)
}

// parseLine returns the command and branchname in the given line
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
