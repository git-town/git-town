package steps

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/step"
)

// JSON is used to store a step in JSON.
type JSON struct { //nolint:musttag // JSONStep uses a custom serialization algorithm
	Step step.Step
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

func DetermineStep(stepType string) step.Step { //nolint:ireturn
	switch stepType {
	case "AbortMerge":
		return &step.AbortMerge{}
	case "AbortRebase":
		return &step.AbortRebase{}
	case "AddToPerennialBranches":
		return &step.AddToPerennialBranches{}
	case "ChangeParent":
		return &step.ChangeParent{}
	case "Checkout":
		return &step.Checkout{}
	case "CheckoutIfExists":
		return &step.CheckoutIfExists{}
	case "CommitOpenChanges":
		return &step.CommitOpenChanges{}
	case "ConnectorMergeProposal":
		return &step.ConnectorMergeProposal{}
	case "ContinueMerge":
		return &step.ContinueMerge{}
	case "ContinueRebase":
		return &step.ContinueRebase{}
	case "CreateBranch":
		return &step.CreateBranch{}
	case "CreateBranchExistingParent":
		return &step.CreateBranchExistingParent{}
	case "CreateProposal":
		return &step.CreateProposal{}
	case "CreateRemoteBranch":
		return &step.CreateRemoteBranch{}
	case "CreateTrackingBranch":
		return &step.CreateTrackingBranch{}
	case "DeleteLocalBranch":
		return &step.DeleteLocalBranch{}
	case "DeleteParentBranch":
		return &step.DeleteParentBranch{}
	case "DeleteRemoteBranch":
		return &step.DeleteRemoteBranch{}
	case "DeleteTrackingBranch":
		return &step.DeleteTrackingBranch{}
	case "DiscardOpenChanges":
		return &step.DiscardOpenChanges{}
	case "Empty":
		return &step.Empty{}
	case "EnsureHasShippableChanges":
		return &step.EnsureHasShippableChanges{}
	case "FetchUpstream":
		return &step.FetchUpstream{}
	case "ForcePushCurrentBranch":
		return &step.ForcePushCurrentBranch{}
	case "Merge":
		return &step.Merge{}
	case "MergeParent":
		return &step.MergeParent{}
	case "PreserveCheckoutHistory":
		return &step.PreserveCheckoutHistory{}
	case "PullCurrentBranch":
		return &step.PullCurrentBranch{}
	case "PushCurrentBranch":
		return &step.PushCurrentBranch{}
	case "PushTags":
		return &step.PushTags{}
	case "RebaseBranch":
		return &step.RebaseBranch{}
	case "RebaseParent":
		return &step.RebaseParent{}
	case "RemoveFromPerennialBranches":
		return &step.RemoveFromPerennialBranches{}
	case "RemoveGlobalConfig":
		return &step.RemoveGlobalConfig{}
	case "RemoveLocalConfig":
		return &step.RemoveLocalConfig{}
	case "ResetCurrentBranchToSHA":
		return &step.ResetCurrentBranchToSHA{}
	case "ResetRemoteBranchToSHA":
		return &step.ResetRemoteBranchToSHA{}
	case "RestoreOpenChanges":
		return &step.RestoreOpenChanges{}
	case "RevertCommit":
		return &step.RevertCommit{}
	case "SetExistingParent":
		return &step.SetExistingParent{}
	case "SetGlobalConfig":
		return &step.SetGlobalConfig{}
	case "SetLocalConfig":
		return &step.SetLocalConfig{}
	case "SetParent":
		return &step.SetParent{}
	case "SquashMerge":
		return &step.SquashMerge{}
	case "SkipCurrentBranch":
		return &step.SkipCurrentBranch{}
	case "StashOpenChanges":
		return &step.StashOpenChanges{}
	case "UndoLastCommit":
		return &step.UndoLastCommit{}
	case "UpdateProposalTarget":
		return &step.UpdateProposalTarget{}
	}
	return nil
}
