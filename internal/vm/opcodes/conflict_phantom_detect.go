package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type ConflictPhantomDetect struct {
	ParentBranch            Option[gitdomain.LocalBranchName]
	ParentSHA               Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomDetect) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *ConflictPhantomDetect) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *ConflictPhantomDetect) Run(args shared.RunArgs) error {
	quickInfos, err := args.Git.FileConflictQuickInfos(args.Backend)
	if err != nil {
		return err
	}
	fullInfos, err := args.Git.FileConflictFullInfos(args.Backend, quickInfos, self.ParentSHA, args.Config.Value.ValidatedConfigData.MainBranch)
	if err != nil {
		return err
	}
	phantomMergeConflicts := git.DetectPhantomMergeConflicts(fullInfos, self.ParentBranch, args.Config.Value.ValidatedConfigData.MainBranch)
	if err != nil {
		return err
	}
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
