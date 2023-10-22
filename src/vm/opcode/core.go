// Package opcode defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcode

import (
	"errors"

	"github.com/git-town/git-town/v9/src/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (self *undeclaredOpcodeMethods) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (self *undeclaredOpcodeMethods) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (self *undeclaredOpcodeMethods) ShouldAutomaticallyAbortOnError() bool {
	return false
}

func Lookup(opcodeType string) shared.Opcode { //nolint:ireturn
	switch opcodeType {
	case "AbortMerge":
		return &AbortMerge{}
	case "AbortRebase":
		return &AbortRebase{}
	case "AddToPerennialBranches":
		return &AddToPerennialBranches{}
	case "ChangeParent":
		return &ChangeParent{}
	case "Checkout":
		return &Checkout{}
	case "CheckoutIfExists":
		return &CheckoutIfExists{}
	case "CommitOpenChanges":
		return &CommitOpenChanges{}
	case "ConnectorMergeProposal":
		return &ConnectorMergeProposal{}
	case "ContinueMerge":
		return &ContinueMerge{}
	case "ContinueRebase":
		return &ContinueRebase{}
	case "CreateBranch":
		return &CreateBranch{}
	case "CreateBranchExistingParent":
		return &CreateBranchExistingParent{}
	case "CreateProposal":
		return &CreateProposal{}
	case "CreateRemoteBranch":
		return &CreateRemoteBranch{}
	case "CreateTrackingBranch":
		return &CreateTrackingBranch{}
	case "DeleteLocalBranch":
		return &DeleteLocalBranch{}
	case "DeleteParentBranch":
		return &DeleteParentBranch{}
	case "DeleteRemoteBranch":
		return &DeleteRemoteBranch{}
	case "DeleteTrackingBranch":
		return &DeleteTrackingBranch{}
	case "DiscardOpenChanges":
		return &DiscardOpenChanges{}
	case "EnsureHasShippableChanges":
		return &EnsureHasShippableChanges{}
	case "FetchUpstream":
		return &FetchUpstream{}
	case "ForcePushCurrentBranch":
		return &ForcePushCurrentBranch{}
	case "IfElse":
		return &IfElse{}
	case "Merge":
		return &Merge{}
	case "MergeParent":
		return &MergeParent{}
	case "PreserveCheckoutHistory":
		return &PreserveCheckoutHistory{}
	case "PullCurrentBranch":
		return &PullCurrentBranch{}
	case "PushCurrentBranch":
		return &PushCurrentBranch{}
	case "PushTags":
		return &PushTags{}
	case "RebaseBranch":
		return &RebaseBranch{}
	case "RebaseParent":
		return &RebaseParent{}
	case "RemoveFromPerennialBranches":
		return &RemoveFromPerennialBranches{}
	case "RemoveGlobalConfig":
		return &RemoveGlobalConfig{}
	case "RemoveLocalConfig":
		return &RemoveLocalConfig{}
	case "ResetCurrentBranchToSHA":
		return &ResetCurrentBranchToSHA{}
	case "ResetRemoteBranchToSHA":
		return &ResetRemoteBranchToSHA{}
	case "RestoreOpenChanges":
		return &RestoreOpenChanges{}
	case "RevertCommit":
		return &RevertCommit{}
	case "SetExistingParent":
		return &SetExistingParent{}
	case "SetGlobalConfig":
		return &SetGlobalConfig{}
	case "SetLocalConfig":
		return &SetLocalConfig{}
	case "SetParent":
		return &SetParent{}
	case "SetParentIfBranchExists":
		return &SetParentIfBranchExists{}
	case "SquashMerge":
		return &SquashMerge{}
	case "SkipCurrentBranch":
		return &SkipCurrentBranch{}
	case "StashOpenChanges":
		return &StashOpenChanges{}
	case "UndoLastCommit":
		return &UndoLastCommit{}
	case "UpdateProposalTarget":
		return &UpdateProposalTarget{}
	}
	return nil
}
