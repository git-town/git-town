package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// MergeParentResolvePhantomConflicts merges the given parent branch into the current branch.
type MergeParentResolvePhantomConflicts struct {
	CurrentBranch     gitdomain.LocalBranchName
	CurrentParent     gitdomain.BranchName              // the currently active parent, after all remotely deleted parents were removed
	InitialParentName Option[gitdomain.LocalBranchName] // name of the original parent when Git Town started
	InitialParentSHA  Option[gitdomain.SHA]             // SHA of the original parent when Git Town started
}

func (self *MergeParentResolvePhantomConflicts) Abort() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeParentResolvePhantomConflicts) Continue() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *MergeParentResolvePhantomConflicts) Run(args shared.RunArgs) error {
	err := args.Git.MergeBranchNoEdit(args.Frontend, self.CurrentParent)
	if err == nil || args.Config.Value.NormalConfig.AutoResolve.NoAutoResolve() {
		return err
	}
	args.PrependOpcodes(&ConflictMergePhantomResolveAll{
		CurrentBranch: self.CurrentBranch,
		ParentBranch:  self.InitialParentName,
		ParentSHA:     self.InitialParentSHA,
	})
	return nil
}
