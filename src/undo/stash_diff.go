package undo

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/opcode"
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

func (self StashDiff) Program() opcode.Program {
	result := opcode.Program{}
	for ; self.EntriesAdded > 0; self.EntriesAdded-- {
		result.Add(&opcode.RestoreOpenChanges{})
	}
	return result
}
