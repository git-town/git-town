package steps

import (
	"github.com/git-town/git-town/src/git"
)

// PushTagsStep pushes newly created Git tags to the remote.
type PushTagsStep struct {
	NoOpStep
}

// Run executes this step.
func (step *PushTagsStep) Run(repo *git.ProdRepo) error {
	return repo.Logging.PushTags()
}
