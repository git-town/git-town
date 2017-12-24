package steps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/util"
)

func hasSavedState(command string) bool {
	filename := getRunResultFilename(command)
	return util.DoesFileExist(filename)
}

func clearSavedState(command string) {
	if hasSavedState(command) {
		exit.If(os.Remove(getRunResultFilename(command)))
	}
}

func loadState(command string) RunState {
	var serializedRunState SerializedRunState
	if hasSavedState(command) {
		content, err := ioutil.ReadFile(getRunResultFilename(command))
		exit.If(err)
		err = json.Unmarshal(content, &serializedRunState)
		exit.If(err)
	} else {
		serializedRunState.AbortStep = SerializedStep{Type: "*NoOpStep"}
	}
	return RunState{
		AbortStep:    deserializeStep(serializedRunState.AbortStep),
		Command:      command,
		RunStepList:  deserializeStepList(serializedRunState.RunSteps),
		UndoStepList: deserializeStepList(serializedRunState.UndoSteps),
	}
}

// nolint: gocyclo
func deserializeStep(serializedStep SerializedStep) Step {
	switch serializedStep.Type {
	case "*AbortMergeBranchStep":
		return &AbortMergeBranchStep{}
	case "*AbortRebaseBranchStep":
		return &AbortRebaseBranchStep{}
	case "*AddToPerennialBranches":
		step := AddToPerennialBranches{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*ChangeDirectoryStep":
		step := ChangeDirectoryStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*CheckoutBranchStep":
		step := CheckoutBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*ContinueMergeBranchStep":
		return &ContinueMergeBranchStep{}
	case "*ContinueRebaseBranchStep":
		return &ContinueRebaseBranchStep{}
	case "*CreateBranchStep":
		step := CreateBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*CreateAndCheckoutBranchStep":
		step := CreateAndCheckoutBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*CreatePullRequestStep":
		step := CreatePullRequestStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*CreateRemoteBranchStep":
		step := CreateRemoteBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*CreateTrackingBranchStep":
		step := CreateTrackingBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*DeleteLocalBranchStep":
		step := DeleteLocalBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*DeleteParentBranchStep":
		step := DeleteParentBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*DeleteRemoteBranchStep":
		step := DeleteRemoteBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*DriverMergePullRequestStep":
		step := DriverMergePullRequestStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*EnsureHasShippableChangesStep":
		step := EnsureHasShippableChangesStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*MergeBranchStep":
		step := MergeBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*NoOpStep":
		return &NoOpStep{}
	case "*PreserveCheckoutHistoryStep":
		step := PreserveCheckoutHistoryStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*PullBranchStep":
		step := PullBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*PushBranchAfterCurrentBranchSteps":
		return &PushBranchAfterCurrentBranchSteps{}
	case "*PushBranchStep":
		step := PushBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*PushTagsStep":
		return &PushTagsStep{}
	case "*RebaseBranchStep":
		step := RebaseBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*RemoveFromPerennialBranches":
		step := RemoveFromPerennialBranches{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*ResetToShaStep":
		step := ResetToShaStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*RestoreOpenChangesStep":
		return &RestoreOpenChangesStep{}
	case "*RevertCommitStep":
		step := RevertCommitStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*SetParentBranchStep":
		step := SetParentBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*SquashMergeBranchStep":
		step := SquashMergeBranchStep{}
		exit.If(json.Unmarshal(serializedStep.Data, &step))
		return &step
	case "*SkipCurrentBranchSteps":
		return &SkipCurrentBranchSteps{}
	case "*StashOpenChangesStep":
		return &StashOpenChangesStep{}
	}
	log.Fatal(fmt.Sprintf("Cannot deserialize steps: %s %s", serializedStep.Type, serializedStep.Data))
	return nil
}

func deserializeStepList(serializedSteps []SerializedStep) (stepList StepList) {
	for _, serializedStep := range serializedSteps {
		stepList.Append(deserializeStep(serializedStep))
	}
	return stepList
}
