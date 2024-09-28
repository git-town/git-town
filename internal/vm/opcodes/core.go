// Package opcodes defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/gohacks"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (self *undeclaredOpcodeMethods) AbortProgram() []shared.Opcode {
	return []shared.Opcode(nil)
}

func (self *undeclaredOpcodeMethods) AutomaticUndoError() error {
	return errors.New("")
}

func (self *undeclaredOpcodeMethods) ContinueProgram() []shared.Opcode {
	return []shared.Opcode(nil)
}

func (self *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (self *undeclaredOpcodeMethods) ShouldUndoOnError() bool {
	return false
}

func (self *undeclaredOpcodeMethods) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode(nil)
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
		&AddToContributionBranches{},
		&AddToObservedBranches{},
		&AddToParkedBranches{},
		&AddToPerennialBranches{},
		&AddToPrototypeBranches{},
		&ChangeParent{},
		&Checkout{},
		&CheckoutIfNeeded{},
		&CheckoutFirstExisting{},
		&CheckoutIfExists{},
		&CheckoutParent{},
		&CheckoutUncached{},
		&ChangeParent{},
		&CommitOpenChanges{},
		&CommitNoEdit{},
		&ConnectorMergeProposal{},
		&ContinueMerge{},
		&ContinueRebase{},
		&ContinueRebaseIfNeeded{},
		&CreateAndCheckoutBranchExistingParent{},
		&CreateBranch{},
		&CreateProposal{},
		&CreateRemoteBranch{},
		&CreateTrackingBranch{},
		&DeleteLocalBranch{},
		&DeleteParentBranch{},
		&DeleteTrackingBranch{},
		&DiscardOpenChanges{},
		&DropStash{},
		&EndOfBranchProgram{},
		&EnsureHasShippableChanges{},
		&FetchUpstream{},
		&ForcePush{},
		&ForcePushCurrentBranch{},
		&DeleteBranchIfEmptyAtRuntime{},
		&Merge{},
		&MergeBranchNoEdit{},
		&MergeParent{},
		&PreserveCheckoutHistory{},
		&PullCurrentBranch{},
		&PushCurrentBranch{},
		&PushTags{},
		&RebaseBranch{},
		&RebaseFeatureTrackingBranch{},
		&RebaseParent{},
		&RemoveBranchFromLineage{},
		&RemoveFromContributionBranches{},
		&RemoveFromObservedBranches{},
		&RemoveFromParkedBranches{},
		&RemoveFromPerennialBranches{},
		&RemoveFromPrototypeBranches{},
		&RemoveGlobalConfig{},
		&RemoveLocalConfig{},
		&RenameBranch{},
		&ResetCurrentBranch{},
		&ResetCurrentBranchToParent{},
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
		&StageOpenChanges{},
		&StashOpenChanges{},
		&SquashMerge{},
		&UndoLastCommit{},
		&UpdateProposalBase{},
		&UpdateProposalHead{},
	} //exhaustruct:ignore
}
