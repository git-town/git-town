package undo

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {
	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches domain.BranchInfos

	// the branch that was checked out at the time the snapshot was taken
	Active domain.LocalBranchName
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{
		Branches: domain.BranchInfos{},
		Active:   domain.LocalBranchName{},
	}
}

// TODO: rename to Spans.
func (b BranchesSnapshot) Span(afterSnapshot BranchesSnapshot) BranchSpans {
	result := BranchSpans{}
	for _, before := range b.Branches {
		after := afterSnapshot.Branches.FindMatchingRecord(before)
		result = append(result, BranchSpan{
			Before: before,
			After:  after,
		})
	}
	for _, after := range afterSnapshot.Branches {
		if b.Branches.FindMatchingRecord(after).IsEmpty() {
			result = append(result, BranchSpan{
				Before: domain.EmptyBranchInfo(),
				After:  after,
			})
		}
	}
	return result
}
