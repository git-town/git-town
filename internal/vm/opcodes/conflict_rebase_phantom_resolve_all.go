package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ConflictRebasePhantomResolveAll struct {
	CurrentBranch           gitdomain.LocalBranchName
	BranchToRebaseOnto      gitdomain.BranchName
	ParentSHA               Option[gitdomain.SHA]
	Resolution              gitdomain.ConflictResolution
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictRebasePhantomResolveAll) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *ConflictRebasePhantomResolveAll) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinue{},
	}
}

func (self *ConflictRebasePhantomResolveAll) Run(args shared.RunArgs) error {
	quickInfos, err := args.Git.FileConflictQuickInfos(args.Backend)
	if err != nil {
		return err
	}
	rootBranch := args.Config.Value.NormalConfig.Lineage.Root(self.CurrentBranch)
	fullInfos, err := args.Git.FileConflictFullInfos(args.Backend, quickInfos, self.BranchToRebaseOnto.Location(), rootBranch)
	if err != nil {
		return err
	}
	phantomRebaseConflicts := git.DetectPhantomRebaseConflicts(fullInfos, self.BranchToRebaseOnto, rootBranch)
	newOpcodes := []shared.Opcode{}
	for _, phantomMergeConflict := range phantomRebaseConflicts {
		newOpcodes = append(newOpcodes, &ConflictPhantomResolve{
			FilePath:   phantomMergeConflict.FilePath,
			Resolution: self.Resolution,
		})
	}
	newOpcodes = append(newOpcodes, &ConflictMergePhantomFinalize{})
	args.PrependOpcodes(newOpcodes...)
	return nil
}
