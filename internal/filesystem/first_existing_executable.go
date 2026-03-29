package filesystem

import (
	"os/exec"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FirstExistingExecutable provides the first existing executable from the given list of executable names.
func FirstExistingExecutable(commands []string) Option[string] {
	for _, command := range commands {
		executable, err := exec.LookPath(command)
		if err == nil && len(executable) > 0 {
			return Some(command)
		}
	}
	return None[string]()
}
