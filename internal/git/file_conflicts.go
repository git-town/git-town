package git

import "github.com/git-town/git-town/v22/internal/subshell/subshelldomain"

type FileConflicts []FileConflict

func (self FileConflicts) Debug(querier subshelldomain.Querier) {
	for _, fileConflict := range self {
		fileConflict.Debug(querier)
	}
}
