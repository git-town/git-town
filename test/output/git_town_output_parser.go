package output

import (
	"strings"

	"github.com/acarl005/stripansi"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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

const (
	CommandTypeFrontend = CommandType("frontend")
	CommandTypeBackend  = CommandType("backend")
)

func (self CommandType) String() string { return string(self) }

// GitCommandsInGitTownOutput provides the Git commands mentioned in the given Git Town output.
func GitCommandsInGitTownOutput(output string) []ExecutedGitCommand {
	result := []ExecutedGitCommand{}
	for _, line := range strings.Split(output, "\n") {
		if lineContainsFrontendCommand(line) {
			if line, hasLine := parseFrontendLine(line).Get(); hasLine {
				result = append(result, line)
			}
		} else if lineContainsBackendCommand(line) {
			result = append(result, parseBackendLine(line))
		}
	}
	return result
}

const backendCommandLineBeginning = "(verbose) "

func lineContainsBackendCommand(line string) bool {
	return strings.HasPrefix(line, backendCommandLineBeginning)
}

const frontendCommandLineBeginning = "\x1b[1m" // "\e[1m"

func lineContainsFrontendCommand(line string) bool {
	return strings.HasPrefix(line, frontendCommandLineBeginning)
}

func parseBackendLine(line string) ExecutedGitCommand {
	command := strings.TrimPrefix(line, backendCommandLineBeginning)
	command = stripansi.Strip(command)
	return ExecutedGitCommand{
		Branch:      "",
		Command:     command,
		CommandType: CommandTypeBackend,
	}
}

func parseFrontendLine(line string) Option[ExecutedGitCommand] {
	line = stripansi.Strip(line)
	if line == "" {
		return None[ExecutedGitCommand]()
	}
	// extract branch name if it exists
	branch := ""
	if line[0] == '[' {
		// line contains a branch name
		closingParent := strings.IndexRune(line, ']')
		branch = line[1:closingParent]
		line = line[closingParent+2:]
	}
	return Some(ExecutedGitCommand{
		Branch:      branch,
		Command:     line,
		CommandType: CommandTypeFrontend,
	})
}
