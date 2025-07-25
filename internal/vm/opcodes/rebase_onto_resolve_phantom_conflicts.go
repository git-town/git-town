package opcodes

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RebaseOntoResolvePhantomConflicts struct {
	BranchToRebaseOnto      gitdomain.BranchName
	CommitsToRemove         gitdomain.Location
	CurrentBranch           gitdomain.LocalBranchName
	Resolution              gitdomain.ConflictResolution
	Upstream                Option[gitdomain.LocalBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoResolvePhantomConflicts) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoResolvePhantomConflicts) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoResolvePhantomConflicts) Run(args shared.RunArgs) error {
	// Fix for https://github.com/git-town/git-town/issues/4942.
	// Waiting here in end-to-end tests to ensure new timestamps for the rebased commits,
	// which avoids flaky end-to-end tests.
	if subshell.IsInTest() {
		time.Sleep(1 * time.Second)
	}
	if err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove, self.Upstream); err != nil {
		conflictingFiles, err := args.Git.FileConflictQuickInfos(args.Backend)
		if err != nil {
			return fmt.Errorf("cannot determine conflicting files after rebase: %w", err)
		}
		rootBranch := args.Config.Value.NormalConfig.Lineage.Root(self.CurrentBranch)
		fullInfos, err := args.Git.FileConflictFullInfos(args.Backend, conflictingFiles, self.BranchToRebaseOnto.Location(), rootBranch)
		if err != nil {
			return err
		}
		phantomRebaseConflicts := git.DetectPhantomRebaseConflicts(fullInfos, self.BranchToRebaseOnto, rootBranch)
		newOpcodes := []shared.Opcode{}
		for _, phantomRebaseConflict := range phantomRebaseConflicts {
			newOpcodes = append(newOpcodes, &ConflictPhantomResolve{
				FilePath:   phantomRebaseConflict.FilePath,
				Resolution: self.Resolution,
			})
		}
		newOpcodes = append(newOpcodes, &ConflictPhantomFinalize{})
		args.PrependOpcodes(newOpcodes...)
		if err = args.Git.ContinueRebase(args.Frontend); err != nil {
			return fmt.Errorf("cannot continue rebase: %w", err)
		}
	}
	return nil
}
