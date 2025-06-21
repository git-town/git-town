package gh

import (
	"os/exec"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func Detect() Option[string] {
	// detect gh executable
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return None[string]()
	}
	return Some(ghPath)
}
