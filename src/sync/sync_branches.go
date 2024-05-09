package sync

import (
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(args BranchesProgramArgs) {
	for _, branch := range args.BranchesToSync {
		BranchProgram(branch, args.BranchProgramArgs)
	}
	previousBranch, hasPreviousBranch := args.PreviousBranch.Get()
	finalBranchCandidates := gitdomain.LocalBranchNames{args.InitialBranch}
	if hasPreviousBranch {
		finalBranchCandidates = append(finalBranchCandidates, previousBranch)
	}
	if hasPreviousBranch {
		args.Program.Add(&opcodes.CheckoutFirstExisting{
			Branches:   finalBranchCandidates,
			MainBranch: args.Config.MainBranch,
		})
	}
	if args.Remotes.HasOrigin() && args.ShouldPushTags && args.Config.IsOnline() {
		args.Program.Add(&opcodes.PushTags{})
	}
	previousbranchCandidates := gitdomain.LocalBranchNames{}
	if hasPreviousBranch {
		previousbranchCandidates = append(previousbranchCandidates, previousBranch)
	}
	cmdhelpers.Wrap(args.Program, cmdhelpers.WrapOptions{
		DryRun:           args.DryRun,
		RunInGitRoot:     true,
		StashOpenChanges: args.HasOpenChanges,
		// TODO: only add args.PreviousBranch if it isn't in another workspace
		PreviousBranchCandidates: previousbranchCandidates,
	})
}

type BranchesProgramArgs struct {
	BranchProgramArgs
	BranchesToSync gitdomain.BranchInfos
	DryRun         bool
	HasOpenChanges bool
	InitialBranch  gitdomain.LocalBranchName
	PreviousBranch Option[gitdomain.LocalBranchName]
	ShouldPushTags bool
}
