package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v20/internal/git"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

type ConflictPhantomResolveAll struct {
	ParentBranch            Option[gitdomain.LocalBranchName]
	ParentSHA               Option[gitdomain.SHA]
	Resolution              gitdomain.ConflictResolution
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomResolveAll) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *ConflictPhantomResolveAll) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *ConflictPhantomResolveAll) Run(args shared.RunArgs) error {
	parentSHA, hasParentSHA := self.ParentSHA.Get()
	if !hasParentSHA {
		return errors.New(messages.ConflictMerge)
	}
	quickInfos, err := args.Git.FileConflictQuickInfos(args.Backend)
	if err != nil {
		return err
	}
	mainBranch := args.Config.Value.ValidatedConfigData.MainBranch
	fullInfos, err := args.Git.FileConflictFullInfos(args.Backend, quickInfos, parentSHA.Location(), mainBranch)
	if err != nil {
		return err
	}
	phantomMergeConflicts := git.DetectPhantomMergeConflicts(fullInfos, self.ParentBranch, mainBranch)
	newOpcodes := make([]shared.Opcode, len(phantomMergeConflicts)+1)
	for p, phantomMergeConflict := range phantomMergeConflicts {
		newOpcodes[p] = &ConflictPhantomResolve{
			FilePath:   phantomMergeConflict.FilePath,
			Resolution: self.Resolution,
		}
	}
	newOpcodes[len(phantomMergeConflicts)] = &ConflictPhantomFinalize{}
	args.PrependOpcodes(newOpcodes...)
	return nil
}
