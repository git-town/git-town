package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// MergeParentResolvePhantomConflicts merges the given parent branch into the current branch.
type MergeParentResolvePhantomConflicts struct {
	CurrentParent           gitdomain.BranchName              // the currently active parent, after all remotely deleted parents were removed
	InitialParentName       Option[gitdomain.LocalBranchName] // name of the original parent when Git Town started
	InitialParentSHA        Option[gitdomain.SHA]             // SHA of the original parent when Git Town started
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParentResolvePhantomConflicts) Run(args shared.RunArgs) error {
	if err := args.Git.MergeBranchNoEdit(args.Frontend, self.CurrentParent); err != nil {
		args.PrependOpcodes(&ConflictPhantomResolveAll{
			ParentBranch: self.InitialParentName,
			ParentSHA:    self.InitialParentSHA,
			Resolution:   gitdomain.ConflictResolutionOurs,
		})
	}
	return nil
}
