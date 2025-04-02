// Package opcodes defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v18/internal/gohacks"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (self *undeclaredOpcodeMethods) AbortProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) AutomaticUndoError() error {
	return errors.New("")
}

func (self *undeclaredOpcodeMethods) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (self *undeclaredOpcodeMethods) ShouldUndoOnError() bool {
	return false
}

func (self *undeclaredOpcodeMethods) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{}
}

// All provides all existing opcodes.
// This is used to iterate all opcode types.
func All() []shared.Opcode {
	return []shared.Opcode{
		&BranchCreate{},
		&BranchCreateAndCheckoutExistingParent{},
		&BranchCurrentReset{},
		&BranchCurrentResetToParent{},
		&BranchCurrentResetToSHA{},
		&BranchCurrentResetToSHAIfNeeded{},
		&BranchWithRemoteGoneDeleteIfEmptyAtRuntime{},
		&BranchEnsureShippableChanges{},
		&BranchLocalDelete{},
		&BranchLocalDeleteContent{},
		&BranchLocalRename{},
		&BranchRemoteCreate{},
		&BranchRemoteSetToSHA{},
		&BranchRemoteSetToSHAIfNeeded{},
		&BranchReset{},
		&BranchTrackingCreate{},
		&BranchTrackingDelete{},
		&BranchTypeOverrideSet{},
		&BranchTypeOverrideRemove{},
		&ChangesDiscard{},
		&Checkout{},
		&CheckoutFirstExisting{},
		&CheckoutHistoryPreserve{},
		&CheckoutIfExists{},
		&CheckoutIfNeeded{},
		&CheckoutParentOrMain{},
		&CheckoutUncached{},
		&ChangesStage{},
		&ChangesUnstageAll{},
		&Commit{},
		&CommitAutoUndo{},
		&CommitMessageCommentOut{},
		&CommitRevert{},
		&CommitRevertIfNeeded{},
		&CommitWithMessage{},
		&ConfigRemove{},
		&ConfigSet{},
		&ConflictPhantomDetect{},
		&ConflictPhantomFinalize{},
		&ConflictPhantomResolve{},
		&ConnectorProposalMerge{},
		&FetchUpstream{},
		&PushCurrentBranchForce{},
		&LineageBranchRemove{},
		&LineageParentRemove{},
		&LineageParentSet{},
		&LineageParentSetFirstExisting{},
		&LineageParentSetIfExists{},
		&LineageParentSetToGrandParent{},
		&Merge{},
		&MergeAbort{},
		&MergeAlwaysProgram{},
		&MergeContinue{},
		&MergeParentResolvePhantomConflicts{},
		&MergeParentIfNeeded{},
		&MergeSquashProgram{},
		&MessageQueue{},
		&ProgramEndOfBranch{},
		&RebaseAbort{},
		&RebaseBranch{},
		&RebaseContinue{},
		&RebaseContinueIfNeeded{},
		&RebaseOntoKeepDeleted{},
		&RebaseOntoRemoveDeleted{},
		&RebaseParentIfNeeded{},
		&RebaseTrackingBranch{},
		&ProposalCreate{},
		&ProposalUpdateTarget{},
		&ProposalUpdateTargetToGrandParent{},
		&ProposalUpdateSource{},
		&PullCurrentBranch{},
		&PushCurrentBranch{},
		&PushCurrentBranchForceIfNeeded{},
		&PushCurrentBranchIfLocal{},
		&PushCurrentBranchIfNeeded{},
		&PushTags{},
		&RegisterUndoablePerennialCommit{},
		&SnapshotInitialUpdateLocalSHA{},
		&SnapshotInitialUpdateLocalSHAIfNeeded{},
		&StashDrop{},
		&StashOpenChanges{},
		&StashPop{},
		&StashPopIfNeeded{},
		&UndoLastCommit{},
	} //exhaustruct:ignore
}

func Lookup(opcodeType string) shared.Opcode { //nolint:ireturn
	for _, opcode := range All() {
		if gohacks.TypeName(opcode) == opcodeType {
			return opcode
		}
	}
	return nil
}
