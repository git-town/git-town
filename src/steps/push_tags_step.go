package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PushTagsStep pushes newly created Git tags to origin.
type PushTagsStep struct {
	EmptyStep
}

func (step *PushTagsStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.PushTags()
}
