package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ConflictMergePhantomResolveAll struct {
	CurrentBranch gitdomain.LocalBranchName
	ParentBranch  Option[gitdomain.LocalBranchName]
	ParentSHA     Option[gitdomain.SHA]
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
	fileConflicts, err := args.Git.FileConflicts(args.Backend)
	if err != nil {
		return err
	}
	rootBranch := args.Config.Value.NormalConfig.Lineage.Root(self.CurrentBranch)
	mergeConflits, err := args.Git.MergeConflicts(args.Backend, fileConflicts, parentSHA.Location(), rootBranch)
	if err != nil {
		return err
	}
	phantomMergeConflicts := git.DetectPhantomMergeConflicts(mergeConflits, self.ParentBranch, rootBranch)
	newOpcodes := []shared.Opcode{}
	for _, phantomMergeConflict := range phantomMergeConflicts {
		newOpcodes = append(newOpcodes, &ConflictResolve{
			FilePath:   phantomMergeConflict.FilePath,
			Resolution: phantomMergeConflict.Resolution,
		})
	}
	newOpcodes = append(newOpcodes, &ConflictMergePhantomFinalize{})
	args.PrependOpcodes(newOpcodes...)
	return nil
}
