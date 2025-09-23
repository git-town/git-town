package git

import "github.com/git-town/git-town/v22/internal/subshell/subshelldomain"

type MergeConflicts []MergeConflict

func (self MergeConflicts) Debug(querier subshelldomain.Querier) {
	for _, mergeConflict := range self {
		mergeConflict.Debug(querier)
	}
}
