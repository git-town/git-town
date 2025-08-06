package git

import "github.com/git-town/git-town/v21/internal/subshell/subshelldomain"

type FileConflicts []FileConflict

func (quickInfos FileConflicts) Debug(querier subshelldomain.Querier) {
	for _, quickInfo := range quickInfos {
		quickInfo.Debug(querier)
	}
}
