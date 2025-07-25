package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ConflictMergePhantomResolveAll struct {
	CurrentBranch           gitdomain.LocalBranchName
	ParentBranch            Option[gitdomain.LocalBranchName]
	ParentSHA               Option[gitdomain.SHA]
	Resolution              gitdomain.ConflictResolution
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictMergePhantomResolveAll) Abort() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *ConflictMergePhantomResolveAll) Continue() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *ConflictMergePhantomResolveAll) Run(args shared.RunArgs) error {
	parentSHA, hasParentSHA := self.ParentSHA.Get()
	if !hasParentSHA {
		return errors.New(messages.ConflictMerge)
	}
	quickInfos, err := args.Git.FileConflictQuickInfos(args.Backend)
	if err != nil {
		return err
	}
	rootBranch := args.Config.Value.NormalConfig.Lineage.Root(self.CurrentBranch)
	fullInfos, err := args.Git.FileConflictFullInfos(args.Backend, quickInfos, parentSHA.Location(), rootBranch)
	if err != nil {
		return err
	}
	phantomMergeConflicts := git.DetectPhantomMergeConflicts(fullInfos, self.ParentBranch, rootBranch)
	newOpcodes := []shared.Opcode{}
	for _, phantomMergeConflict := range phantomMergeConflicts {
		newOpcodes = append(newOpcodes, &ConflictMergePhantomResolve{
			FilePath:   phantomMergeConflict.FilePath,
			Resolution: self.Resolution,
		})
	}
	newOpcodes = append(newOpcodes, &ConflictMergePhantomFinalize{})
	args.PrependOpcodes(newOpcodes...)
	return nil
}
