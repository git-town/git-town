package program

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// JSON is used to store an opcode in JSON.
type JSON struct { //nolint:musttag // JSON uses a custom serialization algorithm
	Opcode shared.Opcode
}

// MarshalJSON marshals the opcode to JSON.
func (self *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": self.Opcode,
		"type": gohacks.TypeName(self.Opcode),
	})
}

// UnmarshalJSON unmarshals the opcode from JSON.
func (self *JSON) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	err := json.Unmarshal(b, &mapping)
	if err != nil {
		return err
	}
	var opcodeType string
	err = json.Unmarshal(mapping["type"], &opcodeType)
	if err != nil {
		return err
	}
	self.Opcode = DetermineOpcode(opcodeType)
	if self.Opcode == nil {
		return fmt.Errorf(messages.OpcodeUnknown, opcodeType)
	}
	return json.Unmarshal(mapping["data"], &self.Opcode)
}

func DetermineOpcode(opcodeType string) shared.Opcode { //nolint:ireturn
	switch opcodeType {
	case "AbortMerge":
		return &opcode.AbortMerge{}
	case "AbortRebase":
		return &opcode.AbortRebase{}
	case "AddToPerennialBranches":
		return &opcode.AddToPerennialBranches{}
	case "ChangeParent":
		return &opcode.ChangeParent{}
	case "Checkout":
		return &opcode.Checkout{}
	case "CheckoutIfExists":
		return &opcode.CheckoutIfExists{}
	case "CommitOpenChanges":
		return &opcode.CommitOpenChanges{}
	case "ConnectorMergeProposal":
		return &opcode.ConnectorMergeProposal{}
	case "ContinueMerge":
		return &opcode.ContinueMerge{}
	case "ContinueRebase":
		return &opcode.ContinueRebase{}
	case "CreateBranch":
		return &opcode.CreateBranch{}
	case "CreateBranchExistingParent":
		return &opcode.CreateBranchExistingParent{}
	case "CreateProposal":
		return &opcode.CreateProposal{}
	case "CreateRemoteBranch":
		return &opcode.CreateRemoteBranch{}
	case "CreateTrackingBranch":
		return &opcode.CreateTrackingBranch{}
	case "DeleteLocalBranch":
		return &opcode.DeleteLocalBranch{}
	case "DeleteParentBranch":
		return &opcode.DeleteParentBranch{}
	case "DeleteRemoteBranch":
		return &opcode.DeleteRemoteBranch{}
	case "DeleteTrackingBranch":
		return &opcode.DeleteTrackingBranch{}
	case "DiscardOpenChanges":
		return &opcode.DiscardOpenChanges{}
	case "EnsureHasShippableChanges":
		return &opcode.EnsureHasShippableChanges{}
	case "FetchUpstream":
		return &opcode.FetchUpstream{}
	case "ForcePushCurrentBranch":
		return &opcode.ForcePushCurrentBranch{}
	case "Merge":
		return &opcode.Merge{}
	case "MergeParent":
		return &opcode.MergeParent{}
	case "PreserveCheckoutHistory":
		return &opcode.PreserveCheckoutHistory{}
	case "PullCurrentBranch":
		return &opcode.PullCurrentBranch{}
	case "PushCurrentBranch":
		return &opcode.PushCurrentBranch{}
	case "PushTags":
		return &opcode.PushTags{}
	case "RebaseBranch":
		return &opcode.RebaseBranch{}
	case "RebaseParent":
		return &opcode.RebaseParent{}
	case "RemoveFromPerennialBranches":
		return &opcode.RemoveFromPerennialBranches{}
	case "RemoveGlobalConfig":
		return &opcode.RemoveGlobalConfig{}
	case "RemoveLocalConfig":
		return &opcode.RemoveLocalConfig{}
	case "ResetCurrentBranchToSHA":
		return &opcode.ResetCurrentBranchToSHA{}
	case "ResetRemoteBranchToSHA":
		return &opcode.ResetRemoteBranchToSHA{}
	case "RestoreOpenChanges":
		return &opcode.RestoreOpenChanges{}
	case "RevertCommit":
		return &opcode.RevertCommit{}
	case "SetExistingParent":
		return &opcode.SetExistingParent{}
	case "SetGlobalConfig":
		return &opcode.SetGlobalConfig{}
	case "SetLocalConfig":
		return &opcode.SetLocalConfig{}
	case "SetParent":
		return &opcode.SetParent{}
	case "SetParentIfBranchExists":
		return &opcode.SetParentIfBranchExists{}
	case "SquashMerge":
		return &opcode.SquashMerge{}
	case "SkipCurrentBranch":
		return &opcode.SkipCurrentBranch{}
	case "StashOpenChanges":
		return &opcode.StashOpenChanges{}
	case "UndoLastCommit":
		return &opcode.UndoLastCommit{}
	case "UpdateProposalTarget":
		return &opcode.UpdateProposalTarget{}
	}
	return nil
}
