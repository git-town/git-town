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
	case "*AbortMergeBranchStep":
		return &steps.AbortMergeBranchStep{}
	case "*AbortRebaseBranchStep":
		return &steps.AbortRebaseBranchStep{}
	case "*AddToPerennialBranchesStep":
		return &steps.AddToPerennialBranchesStep{}
	case "*CheckoutBranchStep":
		return &steps.CheckoutBranchStep{}
	case "*ContinueMergeBranchStep":
		return &steps.ContinueMergeBranchStep{}
	case "*ContinueRebaseBranchStep":
		return &steps.ContinueRebaseBranchStep{}
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
	case "*ConnectorMergeProposalStep":
		return &steps.ConnectorMergeProposalStep{}
	case "*EnsureHasShippableChangesStep":
		return &steps.EnsureHasShippableChangesStep{}
	case "*FetchUpstreamStep":
		return &steps.FetchUpstreamStep{}
	case "*MergeBranchStep":
		return &steps.MergeBranchStep{}
	case "*NoOpStep":
		return &steps.NoOpStep{}
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
	case "*SquashMergeBranchStep":
		return &steps.SquashMergeBranchStep{}
	case "*SkipCurrentBranchSteps":
		return &steps.SkipCurrentBranchSteps{}
	case "*StashOpenChangesStep":
		return &steps.StashOpenChangesStep{}
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
