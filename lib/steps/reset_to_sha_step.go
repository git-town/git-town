package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	Hard bool
	Sha  string
}

// CreateAbortStep returns the abort step for this step.
func (step ResetToShaStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step ResetToShaStep) CreateContinueStep() Step {
	return NoOpStep{}
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
