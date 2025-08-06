package git

import "github.com/git-town/git-town/v21/internal/subshell/subshelldomain"

type MergeConflicts []MergeConflict

func (mergeConflicts MergeConflicts) Debug(querier subshelldomain.Querier) {
	for _, mergeConflict := range mergeConflicts {
		mergeConflict.Debug(querier)
	}
}
