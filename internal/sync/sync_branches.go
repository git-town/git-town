package sync

import (
	"github.com/git-town/git-town/v14/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/opcodes"
	. "github.com/git-town/git-town/v14/pkg/prelude"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(args BranchesProgramArgs) {
	for _, branch := range args.BranchesToSync {
		BranchProgram(branch, args.BranchProgramArgs)
	}
	previousbranchCandidates := gitdomain.LocalBranchNames{}
	finalBranchCandidates := gitdomain.LocalBranchNames{args.InitialBranch}
	if previousBranch, hasPreviousBranch := args.PreviousBranch.Get(); hasPreviousBranch {
		finalBranchCandidates = append(finalBranchCandidates, previousBranch)
		previousbranchCandidates = append(previousbranchCandidates, previousBranch)
	}
	args.Program.Value.Add(&opcodes.CheckoutFirstExisting{
		Branches:   finalBranchCandidates,
		MainBranch: args.Config.MainBranch,
	})
	if args.Remotes.HasOrigin() && args.ShouldPushTags && args.Config.IsOnline() {
		args.Program.Value.Add(&opcodes.PushTags{})
	}
	cmdhelpers.Wrap(args.Program, cmdhelpers.WrapOptions{
		DryRun:                   args.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         args.HasOpenChanges,
		PreviousBranchCandidates: previousbranchCandidates,
	})
}

type BranchesProgramArgs struct {
	BranchProgramArgs
	BranchesToSync gitdomain.BranchInfos
	DryRun         configdomain.DryRun
	HasOpenChanges bool
	InitialBranch  gitdomain.LocalBranchName
	PreviousBranch Option[gitdomain.LocalBranchName]
	ShouldPushTags bool
}
