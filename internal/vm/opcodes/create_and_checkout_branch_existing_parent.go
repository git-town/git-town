package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// CreateAndCheckoutBranchExistingParent creates a new branch with the first existing entry from the given ancestor list as its parent.
type CreateAndCheckoutBranchExistingParent struct {
	Ancestors               gitdomain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CreateAndCheckoutBranchExistingParent) Run(args shared.RunArgs) error {
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	var ancestorToUse gitdomain.BranchName
	nearestAncestor, hasNearestAncestor := args.Git.FirstExistingBranch(args.Backend, self.Ancestors...).Get()
	if hasNearestAncestor {
		ancestorToUse = nearestAncestor.BranchName()
	} else {
		ancestorToUse = args.Config.Config.MainBranch.AtRemote(gitdomain.RemoteOrigin).BranchName()
	}
	if ancestorToUse == currentBranch.BranchName() {
		return args.Git.CreateAndCheckoutBranch(args.Frontend, self.Branch)
	}
	return args.Git.CreateAndCheckoutBranchWithParent(args.Frontend, self.Branch, ancestorToUse.Location())
}
