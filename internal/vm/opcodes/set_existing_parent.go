package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// SetExistingParent sets the first existing entry in the given ancestor list as the parent branch of the given branch.
type SetExistingParent struct {
	Ancestors               gitdomain.LocalBranchNames
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor, hasNearestAncestor := args.Git.FirstExistingBranch(args.Backend, self.Ancestors...).Get()
	if !hasNearestAncestor {
		nearestAncestor = args.Config.Config.MainBranch
	}
	return args.Config.SetParent(self.Branch, nearestAncestor)
}
