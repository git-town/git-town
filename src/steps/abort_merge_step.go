package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// AbortMergeStep aborts the current merge conflict.
type AbortMergeStep struct {
	EmptyStep
}

func (step *AbortMergeStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.AbortMerge()
}
