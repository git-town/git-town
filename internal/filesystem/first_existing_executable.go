package filesystem

import (
	"os/exec"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FirstExistingExecutable provides the first existing executable from the given list of executable names.
func FirstExistingExecutable(commands []string) Option[string] {
	for _, command := range commands {
		if command == "" {
			continue
		}
		_, err := exec.LookPath(command)
		if err == nil {
			return Some(command)
		}
	}
	return None[string]()
}
