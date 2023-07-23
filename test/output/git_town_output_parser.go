package output

import (
	"strings"
)

// ExecutedGitCommand describes a Git command that was executed by Git Town during testing.
type ExecutedGitCommand struct {
	// Branch contains the branch in which this command ran.
	Branch string

	// Command contains the command executed.
	Command string
}

// GitCommandsInGitTownOutput provides the Git commands mentioned in the given Git Town output.
func GitCommandsInGitTownOutput(output string) []ExecutedGitCommand {
	result := []ExecutedGitCommand{}
	for _, line := range strings.Split(output, "\n") {
		if lineContainsGitTownCommand(line) {
			result = append(result, parseLine(line))
		}
	}
	return result
}

// gitCommandLineBeginning contains the first few characters of lines containing Git commands in Git Town output.
const gitCommandLineBeginning = "\x1b[1m" // "\e[1m"

// lineContainsGitTownCommand indicates whether the given line contains a Git Town command.
func lineContainsGitTownCommand(line string) bool {
	return strings.HasPrefix(line, gitCommandLineBeginning)
}

// parseLine provides the Git Town command and branch name in the given line.
func parseLine(line string) ExecutedGitCommand {
	// NOTE: implementing this without regex because the regex has gotten very complex and hard to maintain
	// remove the color codes at the beginning
	line = strings.TrimPrefix(line, gitCommandLineBeginning)
	// extract branch name if it exists
	branch := ""
	if line[0] == '[' {
		// line contains a branch name
		closingParent := strings.IndexRune(line, ']')
		branch = line[1:closingParent]
		line = line[closingParent+2:]
	}
	return ExecutedGitCommand{Command: strings.TrimSpace(line), Branch: branch}
}
