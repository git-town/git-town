package steps

import (
	"github.com/git-town/git-town/src/git"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

// Run executes this step.
func (step *ResetToShaStep) Run(repo *git.ProdRepo) (err error) {
	currentSha, err := repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	if step.Sha == currentSha {
		return nil
	}
	return repo.Logging.ResetToSha(step.Sha, step.Hard)
}
