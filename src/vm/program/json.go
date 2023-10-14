package program

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/opcode"
)

// JSON is used to store a step in JSON.
type JSON struct { //nolint:musttag // JSONStep uses a custom serialization algorithm
	Step opcode.Opcode
}

// MarshalJSON marshals the step to JSON.
func (js *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": js.Step,
		"type": gohacks.TypeName(js.Step),
	})
}

// UnmarshalJSON unmarshals the step from JSON.
func (js *JSON) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	err := json.Unmarshal(b, &mapping)
	if err != nil {
		return err
	}
	var stepType string
	err = json.Unmarshal(mapping["type"], &stepType)
	if err != nil {
		return err
	}
	js.Step = DetermineStep(stepType)
	if js.Step == nil {
		return fmt.Errorf(messages.RunstateStepUnknown, stepType)
	}
	return json.Unmarshal(mapping["data"], &js.Step)
}

func DetermineStep(stepType string) opcode.Opcode { //nolint:ireturn
	switch stepType {
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
	case "Empty":
		return &opcode.Empty{}
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
