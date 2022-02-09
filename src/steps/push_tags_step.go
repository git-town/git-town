package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// PushTagsStep pushes newly created Git tags to origin.
type PushTagsStep struct {
	NoOpStep
}

func (step *PushTagsStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.PushTags()
}
