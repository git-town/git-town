package steps

import (
	"encoding/json"
	"reflect"
)

// JSONStep is used to store a step in JSON
type JSONStep struct {
	Step Step
}

// MarshalJSON marshals the step list to JSON
func (j *JSONStep) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]interface{}{
		"data": j.Step,
		"type": getTypeName(j.Step),
	})
}

// UnmarshalJSON unmarshals the step list from JSON
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
	switch stepType {
	case "*AbortMergeBranchStep":
		j.Step = &AbortMergeBranchStep{}
	case "*AbortRebaseBranchStep":
		j.Step = &AbortRebaseBranchStep{}
	case "*AddToPerennialBranches":
		j.Step = &AddToPerennialBranches{}
	case "*ChangeDirectoryStep":
		j.Step = &ChangeDirectoryStep{}
	case "*CheckoutBranchStep":
		j.Step = &CheckoutBranchStep{}
	case "*ContinueMergeBranchStep":
		j.Step = &ContinueMergeBranchStep{}
	case "*ContinueRebaseBranchStep":
		j.Step = &ContinueRebaseBranchStep{}
	case "*CreateBranchStep":
		j.Step = &CreateBranchStep{}
	case "*CreateAndCheckoutBranchStep":
		j.Step = &CreateAndCheckoutBranchStep{}
	case "*CreatePullRequestStep":
		j.Step = &CreatePullRequestStep{}
	case "*CreateRemoteBranchStep":
		j.Step = &CreateRemoteBranchStep{}
	case "*CreateTrackingBranchStep":
		j.Step = &CreateTrackingBranchStep{}
	case "*DeleteLocalBranchStep":
		j.Step = &DeleteLocalBranchStep{}
	case "*DeleteParentBranchStep":
		j.Step = &DeleteParentBranchStep{}
	case "*DeleteRemoteBranchStep":
		j.Step = &DeleteRemoteBranchStep{}
	case "*DriverMergePullRequestStep":
		j.Step = &DriverMergePullRequestStep{}
	case "*EnsureHasShippableChangesStep":
		j.Step = &EnsureHasShippableChangesStep{}
	case "*MergeBranchStep":
		j.Step = &MergeBranchStep{}
	case "*NoOpStep":
		j.Step = &NoOpStep{}
	case "*PreserveCheckoutHistoryStep":
		j.Step = &PreserveCheckoutHistoryStep{}
	case "*PullBranchStep":
		j.Step = &PullBranchStep{}
	case "*PushBranchAfterCurrentBranchSteps":
		j.Step = &PushBranchAfterCurrentBranchSteps{}
	case "*PushBranchStep":
		j.Step = &PushBranchStep{}
	case "*PushTagsStep":
		j.Step = &PushTagsStep{}
	case "*RebaseBranchStep":
		j.Step = &RebaseBranchStep{}
	case "*RemoveFromPerennialBranches":
		j.Step = &RemoveFromPerennialBranches{}
	case "*ResetToShaStep":
		j.Step = &ResetToShaStep{}
	case "*RestoreOpenChangesStep":
		j.Step = &RestoreOpenChangesStep{}
	case "*RevertCommitStep":
		j.Step = &RevertCommitStep{}
	case "*SetParentBranchStep":
		j.Step = &SetParentBranchStep{}
	case "*SquashMergeBranchStep":
		j.Step = &SquashMergeBranchStep{}
	case "*SkipCurrentBranchSteps":
		j.Step = &SkipCurrentBranchSteps{}
	case "*StashOpenChangesStep":
		j.Step = &StashOpenChangesStep{}
	}
	return json.Unmarshal(*mapping["data"], &j.Step)
}

func getTypeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}
