package opcodes

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// MergeParent merges the given parent branch into the current branch.
type MergeParent struct {
	CurrentBranch           gitdomain.LocalBranchName
	CurrentParent           gitdomain.BranchName              // the currently active parent, after all remotely deleted parents were removed
	InitialParentName       Option[gitdomain.LocalBranchName] // name of the original parent when Git Town started
	InitialParentSHA        Option[gitdomain.SHA]             // SHA of the original parent when Git Town started
	NoAutoResolve           configdomain.AutoResolve
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParent) Run(args shared.RunArgs) error {
	err := args.Git.MergeBranchNoEdit(args.Frontend, self.CurrentParent)
	if err == nil || self.NoAutoResolve {
		return err
	}
	args.PrependOpcodes(&ConflictMergePhantomResolveAll{
		CurrentBranch: self.CurrentBranch,
		ParentBranch:  self.InitialParentName,
		ParentSHA:     self.InitialParentSHA,
	})
	return nil
}
