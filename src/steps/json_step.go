package steps

import (
	"encoding/json"
	"log"
	"reflect"
)

// JSONStep is used to store a step in JSON
type JSONStep struct {
	Step Step
}

// MarshalJSON marshals the step to JSON
func (j *JSONStep) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]interface{}{
		"data": j.Step,
		"type": getTypeName(j.Step),
	})
}

// UnmarshalJSON unmarshals the step from JSON
// nolint: gocyclo
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
	j.Step = getStep(stepType)
	return json.Unmarshal(*mapping["data"], &j.Step)
}

func getStep(stepType string) Step {
	switch stepType {
	case "*AbortMergeBranchStep":
		return &AbortMergeBranchStep{}
	case "*AbortRebaseBranchStep":
		return &AbortRebaseBranchStep{}
	case "*AddToPerennialBranches":
		return &AddToPerennialBranches{}
	case "*ChangeDirectoryStep":
		return &ChangeDirectoryStep{}
	case "*CheckoutBranchStep":
		return &CheckoutBranchStep{}
	case "*ContinueMergeBranchStep":
		return &ContinueMergeBranchStep{}
	case "*ContinueRebaseBranchStep":
		return &ContinueRebaseBranchStep{}
	case "*CreateBranchStep":
		return &CreateBranchStep{}
	case "*CreateAndCheckoutBranchStep":
		return &CreateAndCheckoutBranchStep{}
	case "*CreatePullRequestStep":
		return &CreatePullRequestStep{}
	case "*CreateRemoteBranchStep":
		return &CreateRemoteBranchStep{}
	case "*CreateTrackingBranchStep":
		return &CreateTrackingBranchStep{}
	case "*DeleteLocalBranchStep":
		return &DeleteLocalBranchStep{}
	case "*DeleteParentBranchStep":
		return &DeleteParentBranchStep{}
	case "*DeleteRemoteBranchStep":
		return &DeleteRemoteBranchStep{}
	case "*DriverMergePullRequestStep":
		return &DriverMergePullRequestStep{}
	case "*EnsureHasShippableChangesStep":
		return &EnsureHasShippableChangesStep{}
	case "*MergeBranchStep":
		return &MergeBranchStep{}
	case "*NoOpStep":
		return &NoOpStep{}
	case "*PreserveCheckoutHistoryStep":
		return &PreserveCheckoutHistoryStep{}
	case "*PullBranchStep":
		return &PullBranchStep{}
	case "*PushBranchAfterCurrentBranchSteps":
		return &PushBranchAfterCurrentBranchSteps{}
	case "*PushBranchStep":
		return &PushBranchStep{}
	case "*PushTagsStep":
		return &PushTagsStep{}
	case "*RebaseBranchStep":
		return &RebaseBranchStep{}
	case "*RemoveFromPerennialBranches":
		return &RemoveFromPerennialBranches{}
	case "*ResetToShaStep":
		return &ResetToShaStep{}
	case "*RestoreOpenChangesStep":
		return &RestoreOpenChangesStep{}
	case "*RevertCommitStep":
		return &RevertCommitStep{}
	case "*SetParentBranchStep":
		return &SetParentBranchStep{}
	case "*SquashMergeBranchStep":
		return &SquashMergeBranchStep{}
	case "*SkipCurrentBranchSteps":
		return &SkipCurrentBranchSteps{}
	case "*StashOpenChangesStep":
		return &StashOpenChangesStep{}
	default:
		log.Fatalf("Unknown step type: %s", stepType)
		return nil
	}
}

func getTypeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}
