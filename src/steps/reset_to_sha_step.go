package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

func (step *ResetToShaStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) (err error) {
	currentSha, err := repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	if step.Sha == currentSha {
		return nil
	}
	return repo.Logging.ResetToSha(step.Sha, step.Hard)
}
