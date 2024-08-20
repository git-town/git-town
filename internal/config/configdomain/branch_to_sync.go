package configdomain

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

// all the information needed to sync a branch
type BranchToSync struct {
	BranchInfo         gitdomain.BranchInfo
	FirstCommitMessage Option[gitdomain.CommitMessage] // commit message of the first commit on this branch
}
