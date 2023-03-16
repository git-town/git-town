package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// AbortMergeStep aborts the current merge conflict.
type AbortMergeStep struct {
	EmptyStep
}

func (step *AbortMergeStep) Run(repo *git.PublicRepo, connector hosting.Connector) error {
	return repo.AbortMerge()
}
