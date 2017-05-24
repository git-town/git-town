package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

// Run executes this step.
func (step ResetToShaStep) Run() error {
	if step.Sha == git.GetCurrentSha() {
		return nil
	}
	cmd := []string{"git", "reset"}
	if step.Hard {
		cmd = append(cmd, "--hard")
	}
	cmd = append(cmd, step.Sha)
	return script.RunCommand(cmd...)
}
