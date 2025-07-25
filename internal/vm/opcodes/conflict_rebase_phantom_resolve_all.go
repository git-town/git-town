package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConflictRebasePhantomResolveAll struct {
	CurrentBranch           gitdomain.LocalBranchName
	BranchToRebaseOnto      gitdomain.BranchName
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
	phantomConflicts := git.DetectPhantomRebaseConflicts(fullInfos, self.BranchToRebaseOnto, rootBranch)
	newOpcodes := []shared.Opcode{}
	for _, phantomConflict := range phantomConflicts {
		newOpcodes = append(newOpcodes, &ConflictPhantomResolve{
			FilePath:   phantomConflict.FilePath,
			Resolution: gitdomain.ConflictResolutionTheirs,
		})
	}
	newOpcodes = append(newOpcodes, &ConflictRebasePhantomFinalize{})
	args.PrependOpcodes(newOpcodes...)
	return nil
}
