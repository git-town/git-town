package configdomain

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type BranchesToSync []BranchToSync

func (self BranchesToSync) FindByBranch(needle gitdomain.LocalBranchName) Option[BranchToSync] {
	for _, branchToSync := range self {
		if branchName, hasBranchName := branchToSync.BranchInfo.LocalName.Get(); hasBranchName {
			if branchName == needle {
				return Some(branchToSync)
			}
		}
	}
	return None[BranchToSync]()
}
