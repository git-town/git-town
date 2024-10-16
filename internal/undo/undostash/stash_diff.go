package undostash

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
)

// StashDiff describes the changes made to the Git stash.
type StashDiff struct {
	// the number of entries added to the Git stash (positive = entries added, negative = entries removed, 0 = nothing added)
	EntriesAdded int
}

func NewStashDiff(before, after gitdomain.StashSize) StashDiff {
	return StashDiff{
		EntriesAdded: int(after) - int(before),
	}
}

func (self StashDiff) Program() program.Program {
	result := program.Program{}
	for ; self.EntriesAdded > 0; self.EntriesAdded-- {
		result.Add(&opcodes.StashPopIfNeeded{})
	}
	return result
}
