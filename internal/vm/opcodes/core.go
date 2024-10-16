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
		&BranchCreate{},
		&BranchCreateAndCheckoutExistingParent{},
		&BranchCurrentReset{},
		&BranchCurrentResetToParent{},
		&BranchCurrentResetToSHA{},
		&BranchCurrentResetToSHAIfNeeded{},
		&BranchDeleteIfEmptyAtRuntime{},
		&BranchEnsureShippableChanges{},
		&BranchLocalDelete{},
		&BranchLocalRename{},
		&BranchParentDelete{},
		&BranchRemoteCreate{},
		&BranchRemoteSetToSHA{},
		&BranchRemoteSetToSHAIfNeeded{},
		&BranchReset{},
		&BranchTrackingCreate{},
		&BranchTrackingDelete{},
		&BranchesContributionAdd{},
		&BranchesContributionRemove{},
		&BranchesObservedAdd{},
		&BranchesObservedRemove{},
		&BranchesParkedAdd{},
		&BranchesParkedRemove{},
		&BranchesPerennialAdd{},
		&BranchesPerennialRemove{},
		&BranchesPrototypeAdd{},
		&BranchesPrototypeRemove{},
		&ChangesDiscard{},
		&Checkout{},
		&CheckoutFirstExisting{},
		&CheckoutHistoryPreserve{},
		&CheckoutIfExists{},
		&CheckoutIfNeeded{},
		&CheckoutParentOrMain{},
		&CheckoutUncached{},
		&ChangesStage{},
		&Commit{},
		&CommitRevert{},
		&CommitRevertIfNeeded{},
		&CommitAutoUndo{},
		&CommitWithMessage{},
		&ConfigGlobalRemove{},
		&ConfigGlobalSet{},
		&ConfigLocalRemove{},
		&ConfigLocalSet{},
		&ConnectorProposalMerge{},
		&FetchUpstream{},
		&ForcePush{},
		&LineageBranchRemove{},
		&LineageParentRemove{},
		&LineageParentSet{},
		&LineageParentSetFirstExisting{},
		&LineageParentSetIfExists{},
		&LineageParentSetToGrandParent{},
		&Merge{},
		&MergeAbort{},
		&MergeContinue{},
		&MergeParent{},
		&MergeParentIfNeeded{},
		&MergeSquashProgram{},
		&MessageQueue{},
		&ProgramEndOfBranch{},
		&RebaseAbort{},
		&RebaseBranch{},
		&RebaseContinue{},
		&RebaseContinueIfNeeded{},
		&RebaseParentIfNeeded{},
		&RebaseTrackingBranch{},
		&ProposalCreate{},
		&ProposalUpdateBase{},
		&ProposalUpdateBaseToParent{},
		&ProposalUpdateHead{},
		&PullCurrentBranch{},
		&PushCurrentBranch{},
		&PushCurrentBranchForceIfNeeded{},
		&PushCurrentBranchIfLocal{},
		&PushCurrentBranchIfNeeded{},
		&PushTags{},
		&SnapshotInitialUpdateLocalSHA{},
		&SnapshotInitialUpdateLocalSHAIfNeeded{},
		&StashDrop{},
		&StashOpenChanges{},
		&StashPop{},
		&StashPopIfNeeded{},
		&UndoLastCommit{},
	} //exhaustruct:ignore
}
