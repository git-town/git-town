package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// MergeParentResolvePhantomConflicts merges the given parent branch into the current branch.
type MergeParentResolvePhantomConflicts struct {
	CurrentParent           gitdomain.BranchName              // the currently active parent, after all remotely deleted parents were removed
	OriginalParentName      Option[gitdomain.LocalBranchName] // name of the original parent when Git Town started
	OriginalParentSHA       Option[gitdomain.SHA]             // SHA of the original parent when Git Town started
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParentResolvePhantomConflicts) Run(args shared.RunArgs) error {
	err := args.Git.MergeBranchNoEdit(args.Frontend, self.CurrentParent)
	if err != nil {
		args.PrependOpcodes(&ConflictPhantomDetect{
			ParentBranch: self.OriginalParentName,
			ParentSHA:    self.OriginalParentSHA,
		})
	}
	return nil
}
