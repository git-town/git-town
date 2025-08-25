package swap

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type swapGitOperationsProgramArgs struct {
	children    []swapBranch
	current     gitdomain.BranchInfo
	grandParent gitdomain.LocalBranchName
	parent      gitdomain.BranchInfo
	program     Mutable[program.Program]
}

func swapGitOperationsProgram(args swapGitOperationsProgramArgs) {
	args.program.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.grandParent.BranchName(),
			CommitsToRemove:    args.parent.LocalBranchName().Location(),
		},
	)
	args.program.Value.Add(
		&opcodes.Checkout{
			Branch: args.parent.LocalBranchName(),
		},
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.current.LocalBranchName().BranchName(),
			CommitsToRemove:    args.grandParent.Location(),
		},
	)
	if args.current.HasTrackingBranch() {
		args.program.Value.Add(
			&opcodes.Checkout{
				Branch: args.current.LocalBranchName(),
			},
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: args.current.LocalBranchName(), ForceIfIncludes: true},
		)
	}
	if args.parent.HasTrackingBranch() {
		args.program.Value.Add(
			&opcodes.Checkout{
				Branch: args.parent.LocalBranchName(),
			},
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: args.parent.LocalBranchName(), ForceIfIncludes: true},
		)
	}
	for _, child := range args.children {
		args.program.Value.Add(
			&opcodes.Checkout{
				Branch: child.name,
			},
		)
		oldBranchSHA, hasOldBranchSHA := args.current.LocalSHA.Get()
		if !hasOldBranchSHA {
			oldBranchSHA = args.current.RemoteSHA.GetOrDefault()
		}
		args.program.Value.Add(
			&opcodes.RebaseOnto{
				BranchToRebaseOnto: args.parent.LocalBranchName().BranchName(),
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

	args.program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.current.LocalBranchName()})
}
