package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// MergeParent merges the given parent branch into the current branch.
type MergeParent struct {
	CurrentParent           gitdomain.BranchName              // the currently active parent, after all remotely deleted parents were removed
	OriginalParentName      Option[gitdomain.LocalBranchName] // name of the original parent when Git Town started
	OriginalParentSHA       Option[gitdomain.SHA]             // SHA of the original parent when Git Town started
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParent) Run(args shared.RunArgs) error {
	err := args.Git.MergeBranchNoEdit(args.Frontend, self.CurrentParent)
	if err != nil {
		args.PrependOpcodes(&ConflictPhantomDetect{
			ParentBranch: self.OriginalParentName,
			ParentSHA:    self.OriginalParentSHA,
		})
	}
	return nil
}
