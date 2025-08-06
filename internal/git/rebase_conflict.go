package git

import "github.com/git-town/git-town/v21/internal/subshell/subshelldomain"

type RebaseConflict struct{}

func (conflict RebaseConflict) Debug(querier subshelldomain.Querier) string {
	return "TODO"
}

type RebaseConflicts []RebaseConflict

func (conflicts RebaseConflicts) Debug(querier subshelldomain.Querier) string {
	return "TODO"
}
