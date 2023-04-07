package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	EmptyStep
	Hard bool
	Sha  string
}

func (step *ResetToShaStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	currentSha, err := run.Backend.CurrentSha()
	if err != nil {
		return err
	}
	if step.Sha == currentSha {
		return nil
	}
	return run.Frontend.ResetToSha(step.Sha, step.Hard)
}
