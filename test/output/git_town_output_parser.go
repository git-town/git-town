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

	// frontend or backend
	CommandType commandType
}

type commandType string

const commandTypeFrontend = "frontend"
const commandTypeBackend = "backend"

// GitCommandsInGitTownOutput provides the Git commands mentioned in the given Git Town output.
func GitCommandsInGitTownOutput(output string) []ExecutedGitCommand {
	result := []ExecutedGitCommand{}
	for _, line := range strings.Split(output, "\n") {
		if lineContainsFrontendCommand(line) {
			line := parseFrontendLine(line)
			if line != nil {
				result = append(result, *line)
			}
		} else if lineContainsBackendCommand(line) {
			result = append(result, parseBackendLine(line))
		}
	}
	return result
}

// frontendCommandLineBeginning contains the first few characters of lines containing Git commands in Git Town output.
const frontendCommandLineBeginning = "\x1b[1m" // "\e[1m"

// gitCommandLineBeginning contains the first few characters of lines containing Git commands in Git Town output.
const backendCommandLineBeginning = "(debug) " // "\e[1m"

// lineContainsFrontendCommand indicates whether the given line contains a Git Town command.
func lineContainsFrontendCommand(line string) bool {
	return strings.HasPrefix(line, frontendCommandLineBeginning)
}

// lineContainsGitTownCommand indicates whether the given line contains a Git Town command.
func lineContainsBackendCommand(line string) bool {
	return strings.HasPrefix(line, backendCommandLineBeginning)
}

// parseFrontendLine provides the Git Town command and branch name in the given line.
func parseFrontendLine(line string) *ExecutedGitCommand {
	line = strings.TrimPrefix(line, frontendCommandLineBeginning)
	if line == "" {
		return nil
	}
	// extract branch name if it exists
	branch := ""
	if line[0] == '[' {
		// line contains a branch name
		closingParent := strings.IndexRune(line, ']')
		branch = line[1:closingParent]
		line = line[closingParent+2:]
	}
	return &ExecutedGitCommand{Command: line, Branch: branch, CommandType: commandTypeFrontend}
}

// parseLine provides the Git Town command and branch name in the given line.
func parseBackendLine(line string) ExecutedGitCommand {
	command := strings.TrimPrefix(line, backendCommandLineBeginning)
	return ExecutedGitCommand{
		Branch:      "",
		CommandType: commandTypeBackend,
		Command:     command,
	}
}
