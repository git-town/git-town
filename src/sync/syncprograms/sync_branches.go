package syncprograms

import (
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
)

// SyncBranchesProgram syncs all given branches.
func SyncBranchesProgram(args SyncBranchesProgramArgs) {
	for _, branch := range args.BranchesToSync {
		SyncBranchProgram(branch, args.SyncBranchProgramArgs)
	}
	args.Program.Add(&opcode.CheckoutIfExists{Branch: args.InitialBranch})
	if args.Remotes.HasOrigin() && args.ShouldPushTags && args.IsOnline.Bool() {
		args.Program.Add(&opcode.PushTags{})
	}
	cmdhelpers.Wrap(args.Program, cmdhelpers.WrapOptions{
		RunInGitRoot:             true,
		StashOpenChanges:         args.HasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{args.PreviousBranch},
	})
}

type SyncBranchesProgramArgs struct {
	SyncBranchProgramArgs
	BranchesToSync undodomain.BranchInfos
	HasOpenChanges bool
	InitialBranch  gitdomain.LocalBranchName
	PreviousBranch gitdomain.LocalBranchName
	ShouldPushTags bool
}
