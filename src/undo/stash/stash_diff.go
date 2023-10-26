package stash

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/program"
)

// Diff describes the changes made to the Git stash.
type Diff struct {
	// the number of entries added to the Git stash (positive = entries added, negative = entries removed, 0 = nothing added)
	EntriesAdded int
}

func NewDiff(before, after domain.StashSnapshot) Diff {
	return Diff{
		EntriesAdded: int(after) - int(before),
	}
}

func (self Diff) Program() program.Program {
	result := program.Program{}
	for ; self.EntriesAdded > 0; self.EntriesAdded-- {
		result.Add(&opcode.RestoreOpenChanges{})
	}
	return result
}
