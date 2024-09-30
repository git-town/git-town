package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeParentIfNeeded struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParentIfNeeded) Run(args shared.RunArgs) error {
	if parent, hasParent := args.Config.Config.Lineage.Parent(self.CurrentBranch).Get(); hasParent {
		if args.Git.BranchExists(args.Backend, parent) {
			parentActiveIn := true
			args.PrependOpcodes(&MergeParent{
				CurrentBranch:               "",
				ParentActiveInOtherWorktree: false,
				undeclaredOpcodeMethods:     undeclaredOpcodeMethods{},
			})
			// parent is local
		} else {
			// parent isn't local
		}
	}
	return nil
}
