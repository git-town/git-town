package runstate

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/steps"
)

// JSONStep is used to store a step in JSON.
type JSONStep struct { //nolint:musttag // JSONStep uses a custom serialization algorithm
	Step steps.Step
}

// MarshalJSON marshals the step to JSON.
func (j *JSONStep) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": j.Step,
		"type": gohacks.TypeName(j.Step),
	})
}

// UnmarshalJSON unmarshals the step from JSON.
func (j *JSONStep) UnmarshalJSON(b []byte) error {
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
	j.Step = DetermineStep(stepType)
	if j.Step == nil {
		return fmt.Errorf(messages.RunstateStepUnknown, stepType)
	}
	return json.Unmarshal(mapping["data"], &j.Step)
}

func DetermineStep(stepType string) steps.Step {
	switch stepType {
	case "AbortMergeStep":
		return &steps.AbortMergeStep{}
	case "AbortRebaseStep":
		return &steps.AbortRebaseStep{}
	case "AddToPerennialBranchesStep":
		return &steps.AddToPerennialBranchesStep{}
	case "CheckoutStep":
		return &steps.CheckoutStep{}
	case "CheckoutIfExistsStep":
		return &steps.CheckoutIfExistsStep{}
	case "CommitOpenChangesStep":
		return &steps.CommitOpenChangesStep{}
	case "ConnectorMergeProposalStep":
		return &steps.ConnectorMergeProposalStep{}
	case "ContinueMergeStep":
		return &steps.ContinueMergeStep{}
	case "ContinueRebaseStep":
		return &steps.ContinueRebaseStep{}
	case "CreateBranchStep":
		return &steps.CreateBranchStep{}
	case "CreateProposalStep":
		return &steps.CreateProposalStep{}
	case "CreateRemoteBranchStep":
		return &steps.CreateRemoteBranchStep{}
	case "CreateTrackingBranchStep":
		return &steps.CreateTrackingBranchStep{}
	case "DeleteLocalBranchStep":
		return &steps.DeleteLocalBranchStep{}
	case "DeleteParentBranchStep":
		return &steps.DeleteParentBranchStep{}
	case "DeleteRemoteBranchStep":
		return &steps.DeleteRemoteBranchStep{}
	case "DeleteTrackingBranchStep":
		return &steps.DeleteTrackingBranchStep{}
	case "DiscardOpenChangesStep":
		return &steps.DiscardOpenChangesStep{}
	case "EmptyStep":
		return &steps.EmptyStep{}
	case "EnsureHasShippableChangesStep":
		return &steps.EnsureHasShippableChangesStep{}
	case "FetchUpstreamStep":
		return &steps.FetchUpstreamStep{}
	case "ForcePushCurrentBranchStep":
		return &steps.ForcePushCurrentBranchStep{}
	case "MergeStep":
		return &steps.MergeStep{}
	case "PreserveCheckoutHistoryStep":
		return &steps.PreserveCheckoutHistoryStep{}
	case "PullCurrentBranchStep":
		return &steps.PullCurrentBranchStep{}
	case "PushCurrentBranchStep":
		return &steps.PushCurrentBranchStep{}
	case "PushTagsStep":
		return &steps.PushTagsStep{}
	case "RebaseBranchStep":
		return &steps.RebaseBranchStep{}
	case "RemoveFromPerennialBranchesStep":
		return &steps.RemoveFromPerennialBranchesStep{}
	case "RemoveGlobalConfigStep":
		return &steps.RemoveGlobalConfigStep{}
	case "RemoveLocalConfigStep":
		return &steps.RemoveLocalConfigStep{}
	case "ResetCurrentBranchToSHAStep":
		return &steps.ResetCurrentBranchToSHAStep{}
	case "ResetRemoteBranchToSHAStep":
		return &steps.ResetRemoteBranchToSHAStep{}
	case "RestoreOpenChangesStep":
		return &steps.RestoreOpenChangesStep{}
	case "RevertCommitStep":
		return &steps.RevertCommitStep{}
	case "SetGlobalConfigStep":
		return &steps.SetGlobalConfigStep{}
	case "SetLocalConfigStep":
		return &steps.SetLocalConfigStep{}
	case "SetParentStep":
		return &steps.SetParentStep{}
	case "SquashMergeStep":
		return &steps.SquashMergeStep{}
	case "SkipCurrentBranchSteps":
		return &steps.SkipCurrentBranchSteps{}
	case "StashOpenChangesStep":
		return &steps.StashOpenChangesStep{}
	case "UndoLastCommitStep":
		return &steps.UndoLastCommitStep{}
	case "UpdateProposalTargetStep":
		return &steps.UpdateProposalTargetStep{}
	}
	return nil
}
