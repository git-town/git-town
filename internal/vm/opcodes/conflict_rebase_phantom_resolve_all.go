package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConflictRebasePhantomResolveAll struct {
	BranchToRebaseOnto      gitdomain.BranchName
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictRebasePhantomResolveAll) Run(args shared.RunArgs) error {
	fileConflicts, err := args.Git.FileConflicts(args.Backend)
	if err != nil {
		return err
	}
	fileConflicts.Debug(args.Backend)
	rootBranch := args.Config.Value.NormalConfig.Lineage.Root(self.CurrentBranch)
	rebaseConflicts, err := args.Git.RebaseConflicts(args.Backend, fileConflicts, self.BranchToRebaseOnto.Location(), rootBranch)
	if err != nil {
		return err
	}
	rebaseConflicts.Debug(args.Backend)
	phantomConflicts := git.DetectPhantomRebaseConflicts(rebaseConflicts, self.BranchToRebaseOnto, rootBranch)
	newOpcodes := []shared.Opcode{}
	for _, phantomConflict := range phantomConflicts {
		newOpcodes = append(newOpcodes, &ConflictPhantomResolve{
			FilePath:   phantomConflict.FilePath,
			Resolution: phantomConflict.Resolution,
		})
	}
	newOpcodes = append(newOpcodes, &ConflictRebasePhantomFinalize{})
	args.PrependOpcodes(newOpcodes...)
	return nil
}
