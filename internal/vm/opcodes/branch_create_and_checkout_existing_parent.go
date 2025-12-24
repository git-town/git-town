package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchCreateAndCheckoutExistingParent creates a new branch with the first existing entry from the given ancestor list as its parent.
type BranchCreateAndCheckoutExistingParent struct {
	Ancestors gitdomain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	Branch    gitdomain.LocalBranchName
}

func (self *BranchCreateAndCheckoutExistingParent) Run(args shared.RunArgs) error {
	currentBranchOpt, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if !hasCurrentBranch {
		return errors.New(messages.CurrentBranchCannotDetermine)
	}
	var ancestorToUse gitdomain.BranchName
	if nearestAncestor, hasNearestAncestor := args.Git.FirstExistingBranch(args.Backend, self.Ancestors...).Get(); hasNearestAncestor {
		ancestorToUse = nearestAncestor.BranchName()
	} else {
		if mainInfo, hasMainBranch := args.BranchInfos.FindLocalOrRemote(args.Config.Value.ValidatedConfigData.MainBranch).Get(); hasMainBranch {
			ancestorToUse = mainInfo.GetLocalOrRemoteName()
		}
	}
	if ancestorToUse == currentBranch.BranchName() {
		return args.Git.CreateAndCheckoutBranch(args.Frontend, self.Branch)
	}
	return args.Git.CreateAndCheckoutBranchWithParent(args.Frontend, self.Branch, ancestorToUse.Location())
}
