package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConflictPhantomDetect struct {
	ParentBranch            gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomDetect) Run(args shared.RunArgs) error {
	unmergedFiles, err := args.Git.UnmergedFiles(args.Backend)
	if err != nil {
		return err
	}
	phantomMergeConflicts, err := args.Git.DetectPhantomMergeConflicts(args.Backend, unmergedFiles, self.ParentBranch, args.Config.Value.ValidatedConfigData.MainBranch)
	newOpcodes := make([]shared.Opcode, len(phantomMergeConflicts)+1)
	for p, phantomMergeConflict := range phantomMergeConflicts {
		newOpcodes[p] = &ConflictPhantomResolve{
			FilePath: phantomMergeConflict.FilePath,
		}
	}
	newOpcodes[len(phantomMergeConflicts)] = &ConflictPhantomFinalize{}
	args.PrependOpcodes(newOpcodes...)
	return nil
}
