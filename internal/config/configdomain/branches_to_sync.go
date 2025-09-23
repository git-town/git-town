package configdomain

import "github.com/git-town/git-town/v22/internal/git/gitdomain"

type BranchesToSync []BranchToSync

func (self BranchesToSync) BranchNames() gitdomain.LocalBranchNames {
	result := make(gitdomain.LocalBranchNames, len(self))
	for b, branchToSync := range self {
		result[b] = branchToSync.BranchInfo.LocalBranchName()
	}
	return result
}
