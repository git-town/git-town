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
	CommandType CommandType
}

type CommandType string

const CommandTypeFrontend = "frontend"
const CommandTypeBackend = "backend"

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

const frontendCommandLineBeginning = "\x1b[1m" // "\e[1m"

func lineContainsFrontendCommand(line string) bool {
	return strings.HasPrefix(line, frontendCommandLineBeginning)
}

const backendCommandLineBeginning = "(debug) " // "\e[1m"

func lineContainsBackendCommand(line string) bool {
	return strings.HasPrefix(line, backendCommandLineBeginning)
}

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
	return &ExecutedGitCommand{Command: line, Branch: branch, CommandType: CommandTypeFrontend}
}

func parseBackendLine(line string) ExecutedGitCommand {
	command := strings.TrimPrefix(line, backendCommandLineBeginning)
	return ExecutedGitCommand{
		Branch:      "",
		CommandType: CommandTypeBackend,
		Command:     command,
	}
}
