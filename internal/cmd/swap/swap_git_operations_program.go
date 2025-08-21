package swap

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type swapGitOperationsProgramArgs struct {
	branchToSwap        gitdomain.BranchInfo
	childBranchesToSwap []swapBranch
	grandParentBranch   gitdomain.LocalBranchName
	parentBranch        gitdomain.BranchInfo
	program             Mutable[program.Program]
}

func swapGitOperationsProgram(args swapGitOperationsProgramArgs) {
	args.program.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.grandParentBranch.BranchName(),
			CommitsToRemove:    args.parentBranch.LocalBranchName().Location(),
		},
	)

	if args.branchToSwap.HasTrackingBranch() {
		args.program.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: args.branchToSwap.LocalBranchName(), ForceIfIncludes: true},
		)
	}
	args.program.Value.Add(
		&opcodes.Checkout{
			Branch: args.parentBranch.LocalBranchName(),
		},
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.branchToSwap.LocalBranchName().BranchName(),
			CommitsToRemove:    args.grandParentBranch.Location(),
		},
	)
	if args.parentBranch.HasTrackingBranch() {
		args.program.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: args.parentBranch.LocalBranchName(), ForceIfIncludes: true},
		)
	}
	for _, child := range args.childBranchesToSwap {
		args.program.Value.Add(
			&opcodes.Checkout{
				Branch: child.name,
			},
		)
		oldBranchSHA, hasOldBranchSHA := args.branchToSwap.LocalSHA.Get()
		if !hasOldBranchSHA {
			oldBranchSHA = args.branchToSwap.RemoteSHA.GetOrDefault()
		}
		args.program.Value.Add(
			&opcodes.RebaseOnto{
				BranchToRebaseOnto: args.parentBranch.LocalBranchName().BranchName(),
				CommitsToRemove:    oldBranchSHA.Location(),
			},
		)
		if child.info.HasTrackingBranch() {
			args.program.Value.Add(
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   child.name,
					ForceIfIncludes: true,
				},
			)
		}
	}
	args.program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.branchToSwap.LocalBranchName()})
}
