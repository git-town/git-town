// Package opcode defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v15/internal/gohacks"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (self *undeclaredOpcodeMethods) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) CreateAutomaticUndoError() error {
	return errors.New("")
}

func (self *undeclaredOpcodeMethods) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (self *undeclaredOpcodeMethods) ShouldAutomaticallyUndoOnError() bool {
	return false
}

func Lookup(opcodeType string) shared.Opcode { //nolint:ireturn
	for _, opcode := range Types() {
		if gohacks.TypeName(opcode) == opcodeType {
			return opcode
		}
	}
	return nil
}

// Types provides all existing opcodes.
// This is used to iterate all opcode types.
func Types() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
		&AbortRebase{},
		&AddToPerennialBranches{},
		&ChangeParent{},
		&Checkout{},
		&CheckoutFirstExisting{},
		&CheckoutIfExists{},
		&CheckoutParent{},
		&ChangeParent{},
		&CommitOpenChanges{},
		&ConnectorMergeProposal{},
		&ContinueMerge{},
		&ContinueRebase{},
		&CreateAndCheckoutBranchExistingParent{},
		&CreateBranch{},
		&CreateProposal{},
		&CreateRemoteBranch{},
		&CreateTrackingBranch{},
		&DeleteLocalBranch{},
		&DeleteParentBranch{},
		&DeleteTrackingBranch{},
		&DiscardOpenChanges{},
		&EndOfBranchProgram{},
		&EnsureHasShippableChanges{},
		&FetchUpstream{},
		&ForcePushCurrentBranch{},
		&DeleteBranchIfEmptyAtRuntime{},
		&Merge{},
		&MergeParent{},
		&PreserveCheckoutHistory{},
		&PullCurrentBranch{},
		&PushCurrentBranch{},
		&PushTags{},
		&RebaseBranch{},
		&RebaseFeatureTrackingBranch{},
		&RebaseParent{},
		&RemoveBranchFromLineage{},
		&RemoveFromPerennialBranches{},
		&RemoveGlobalConfig{},
		&RemoveLocalConfig{},
		&ResetCurrentBranch{},
		&ResetCurrentBranchToSHA{},
		&ResetRemoteBranchToSHA{},
		&RestoreOpenChanges{},
		&RevertCommit{},
		&SetExistingParent{},
		&SetGlobalConfig{},
		&SetLocalConfig{},
		&SetParent{},
		&SetParentIfBranchExists{},
		&SkipCurrentBranch{},
		&StashOpenChanges{},
		&SquashMerge{},
		&UndoLastCommit{},
		&UpdateProposalTarget{},
	} //exhaustruct:ignore
}
