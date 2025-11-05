package undostash

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
)

// StashDiff describes the changes made to the Git stash.
type StashDiff struct {
	// the number of entries added to the Git stash (positive = entries added, negative = entries removed, 0 = nothing added)
	EntriesAdded int
}

func NewStashDiff(before, after gitdomain.StashSize) StashDiff {
	diff := int(after) - int(before)
	// limit stashes to unpop to at most 1 because Git Town never creates more than 1 stash entry
	diff = min(diff, 1)
	return StashDiff{
		EntriesAdded: diff,
	}
}

func (self StashDiff) Program() program.Program {
	result := program.Program{}
	for ; self.EntriesAdded > 0; self.EntriesAdded-- {
		result.Add(&opcodes.StashPopIfExists{})
	}
	return result
}
