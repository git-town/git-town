package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

// Run executes this step.
func (step *ResetToShaStep) Run() error {
	if step.Sha == git.GetCurrentSha() {
		return nil
	}
	args := []string{"reset"}
	if step.Hard {
		args = append(args, "--hard")
	}
	args = append(args, step.Sha)
	return script.RunCommand("git", args...)
}
