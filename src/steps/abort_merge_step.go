package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// AbortMergeStep aborts the current merge conflict.
type AbortMergeStep struct {
	EmptyStep
}

func (step *AbortMergeStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.AbortMerge()
}
