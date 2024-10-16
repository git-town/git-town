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
		&Checkout{},
		&CheckoutIfNeeded{},
		&CheckoutFirstExisting{},
		&CheckoutIfExists{},
		&CheckoutParent{},
		&CheckoutUncached{},
		&Commit{},
		&CommitWithMessage{},
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
		&LineageSetParent{},
		&Merge{},
		&MergeParent{},
		&MergeParentIfNeeded{},
		&PopStash{},
		&PreserveCheckoutHistory{},
		&PullCurrentBranch{},
		&PushCurrentBranch{},
		&PushCurrentBranchIfLocal{},
		&PushCurrentBranchIfNeeded{},
		&PushTags{},
		&QueueMessage{},
		&RebaseBranch{},
		&RebaseFeatureTrackingBranch{},
		&RebaseParentIfNeeded{},
		&RemoveBranchFromLineage{},
		&RemoveFromContributionBranches{},
		&RemoveFromObservedBranches{},
		&RemoveFromParkedBranches{},
		&RemoveFromPerennialBranches{},
		&RemoveFromPrototypeBranches{},
		&RemoveGlobalConfig{},
		&RemoveLocalConfig{},
		&RemoveParent{},
		&Rename{},
		&ResetBranch{},
		&ResetCurrentBranch{},
		&ResetCurrentBranchToParent{},
		&ResetCurrentBranchToSHA{},
		&ResetCurrentBranchToSHAIfNeeded{},
		&ResetRemoteBranchToSHA{},
		&ResetRemoteBranchToSHAIfNeeded{},
		&RestoreOpenChanges{},
		&RevertCommit{},
		&RevertCommitIfNeeded{},
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
