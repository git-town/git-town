package runstate

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/git-town/git-town/v7/src/steps"
)

// JSONStep is used to store a step in JSON.
type JSONStep struct {
	Step steps.Step
}

// MarshalJSON marshals the step to JSON.
func (j *JSONStep) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": j.Step,
		"type": typeName(j.Step),
	})
}

// UnmarshalJSON unmarshals the step from JSON.
func (j *JSONStep) UnmarshalJSON(b []byte) error {
	var mapping map[string]*json.RawMessage
	err := json.Unmarshal(b, &mapping)
	if err != nil {
		return err
	}
	var stepType string
	err = json.Unmarshal(*mapping["type"], &stepType)
	if err != nil {
		return err
	}
	j.Step = determineStep(stepType)
	return json.Unmarshal(*mapping["data"], &j.Step)
}

func determineStep(stepType string) steps.Step {
	switch stepType {
	case "*AbortMergeStep":
		return &steps.AbortMergeStep{}
	case "*AbortRebaseStep":
		return &steps.AbortRebaseStep{}
	case "*AddToPerennialBranchesStep":
		return &steps.AddToPerennialBranchesStep{}
	case "*CheckoutStep":
		return &steps.CheckoutStep{}
	case "*ConnectorMergeProposalStep":
		return &steps.ConnectorMergeProposalStep{}
	case "*ContinueMergeStep":
		return &steps.ContinueMergeStep{}
	case "*ContinueRebaseStep":
		return &steps.ContinueRebaseStep{}
	case "*CreateBranchStep":
		return &steps.CreateBranchStep{}
	case "*CreateProposalStep":
		return &steps.CreateProposalStep{}
	case "*CreateRemoteBranchStep":
		return &steps.CreateRemoteBranchStep{}
	case "*CreateTrackingBranchStep":
		return &steps.CreateTrackingBranchStep{}
	case "*DeleteLocalBranchStep":
		return &steps.DeleteLocalBranchStep{}
	case "*DeleteOriginBranchStep":
		return &steps.DeleteOriginBranchStep{}
	case "*DeleteParentBranchStep":
		return &steps.DeleteParentBranchStep{}
	case "*DiscardOpenChangesStep":
		return &steps.DiscardOpenChangesStep{}
	case "*EmptyStep":
		return &steps.EmptyStep{}
	case "*EnsureHasShippableChangesStep":
		return &steps.EnsureHasShippableChangesStep{}
	case "*FetchUpstreamStep":
		return &steps.FetchUpstreamStep{}
	case "*MergeStep":
		return &steps.MergeStep{}
	case "*PreserveCheckoutHistoryStep":
		return &steps.PreserveCheckoutHistoryStep{}
	case "*PullBranchStep":
		return &steps.PullBranchStep{}
	case "*PushBranchAfterCurrentBranchSteps":
		return &steps.PushBranchAfterCurrentBranchSteps{}
	case "*PushBranchStep":
		return &steps.PushBranchStep{}
	case "*PushTagsStep":
		return &steps.PushTagsStep{}
	case "*RebaseBranchStep":
		return &steps.RebaseBranchStep{}
	case "*RemoveFromPerennialBranchesStep":
		return &steps.RemoveFromPerennialBranchesStep{}
	case "*ResetToShaStep":
		return &steps.ResetToShaStep{}
	case "*RestoreOpenChangesStep":
		return &steps.RestoreOpenChangesStep{}
	case "*RevertCommitStep":
		return &steps.RevertCommitStep{}
	case "*SetParentStep":
		return &steps.SetParentStep{}
	case "*SquashMergeStep":
		return &steps.SquashMergeStep{}
	case "*SkipCurrentBranchSteps":
		return &steps.SkipCurrentBranchSteps{}
	case "*StashOpenChangesStep":
		return &steps.StashOpenChangesStep{}

	// legacy steps (remove this section in 2026)
	case "*AbortMergeBranchStep":
		return &steps.AbortMergeStep{}
	case "*AbortRebaseBranchStep":
		return &steps.AbortRebaseStep{}
	case "*CheckoutBranchStep":
		return &steps.CheckoutStep{}
	case "*ContinueRebaseBranchStep":
		return &steps.ContinueRebaseStep{}
	case "*MergeBranchStep":
		return &steps.MergeStep{}
	case "*NoOpStep":
		return &steps.EmptyStep{}
	case "*SquashMergeBranchStep":
		return &steps.SquashMergeStep{}

	default:
		log.Fatalf("Unknown step type: %s", stepType)
		return nil
	}
}

func typeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}
