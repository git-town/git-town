package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	EmptyStep
	Hard bool
	Sha  domain.SHA
}

func (step *ResetToShaStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	currentSha, err := run.Backend.CurrentSha()
	if err != nil {
		return err
	}
	if step.Sha == currentSha {
		return nil
	}
	return run.Frontend.ResetToSha(step.Sha, step.Hard)
}
