package undo

import (
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
)

// StashDiff describes the changes made to the Git stash.
type StashDiff struct {
	// the number of entries added to the Git stash (positive = entries added, negative = entries removed, 0 = nothing added)
	EntriesAdded int
}

func (sd StashDiff) Steps() runstate.StepList {
	result := runstate.StepList{}
	if sd.EntriesAdded > 0 {
		for sd.EntriesAdded > 0 {
			result.Append(&steps.RestoreOpenChangesStep{})
		}
	}
	if sd.EntriesAdded < 0 {
		panic("unexpected smaller Git stash found")
	}
	return result
}
