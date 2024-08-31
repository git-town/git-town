package sync

import (
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(args BranchesProgramArgs) {
	for _, branchToSync := range args.BranchesToSync {
		args.BranchProgramArgs.FirstCommitMessage = branchToSync.FirstCommitMessage
		BranchProgram(branchToSync.BranchInfo, args.BranchProgramArgs)
	}
	previousbranchCandidates := []Option[gitdomain.LocalBranchName]{args.PreviousBranch}
	// TODO: make this a list of options
	finalBranchCandidates := gitdomain.LocalBranchNames{args.InitialBranch}
	if previousBranch, hasPreviousBranch := args.PreviousBranch.Get(); hasPreviousBranch {
		finalBranchCandidates = append(finalBranchCandidates, previousBranch)
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
	BranchesToSync []configdomain.BranchToSync
	DryRun         configdomain.DryRun
	HasOpenChanges bool
	InitialBranch  gitdomain.LocalBranchName
	PreviousBranch Option[gitdomain.LocalBranchName]
	ShouldPushTags bool
}
