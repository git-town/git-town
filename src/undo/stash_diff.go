package undo

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/vm/program"
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

func (sd StashDiff) Program() program.Program {
	result := program.Program{}
	for ; sd.EntriesAdded > 0; sd.EntriesAdded-- {
		result.Add(&step.RestoreOpenChanges{})
	}
	return result
}
