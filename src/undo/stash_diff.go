package undo

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
)

// StashDiff describes the changes made to the Git stash.
type StashDiff struct {
	// the number of entries added to the Git stash (positive = entries added, negative = entries removed, 0 = nothing added)
	EntriesAdded int
}

func NewStashDiff(before, after domain.StashSnapshot) StashDiff {
	return StashDiff{
		EntriesAdded: after.Amount - before.Amount,
	}
}

func (sd StashDiff) Steps() runstate.StepList {
	result := runstate.StepList{}
	for ; sd.EntriesAdded > 0; sd.EntriesAdded-- {
		result.Append(&steps.RestoreOpenChangesStep{})
	}
	return result
}
