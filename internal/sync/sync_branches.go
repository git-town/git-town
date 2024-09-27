package sync

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(args BranchesProgramArgs) {
}

type BranchesProgramArgs struct {
	BranchProgramArgs
	BranchesToSync []configdomain.BranchToSync
	DryRun         configdomain.DryRun
	HasOpenChanges bool
	InitialBranch  gitdomain.LocalBranchName
	PreviousBranch Option[gitdomain.LocalBranchName]
	ShouldPushTags bool
}
