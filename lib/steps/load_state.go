package steps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func hasSavedState(command string) bool {
	filename := getRunResultFilename(command)
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func clearSavedState(command string) {
	if hasSavedState(command) {
		os.Remove(getRunResultFilename(command))
	}
}

func loadState(command string) RunState {
	var serializedRunState SerializedRunState
	if hasSavedState(command) {
		content, err := ioutil.ReadFile(getRunResultFilename(command))
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(content, &serializedRunState)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		serializedRunState.AbortStep = SerializedStep{Type: "NoOpStep"}
	}
	return RunState{
		AbortStep:    deserializeStep(serializedRunState.AbortStep),
		Command:      command,
		RunStepList:  deserializeStepList(serializedRunState.RunSteps),
		UndoStepList: deserializeStepList(serializedRunState.UndoSteps),
	}
}

func deserializeStep(serializedStep SerializedStep) Step {
	switch serializedStep.Type {
	case "AbortMergeBranchStep":
		return AbortMergeBranchStep{}
	case "AbortRebaseBranchStep":
		return AbortRebaseBranchStep{}
	case "ChangeDirectoryStep":
		step := ChangeDirectoryStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "CheckoutBranchStep":
		step := CheckoutBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "ContinueMergeBranchStep":
		return ContinueMergeBranchStep{}
	case "ContinueRebaseBranchStep":
		return ContinueRebaseBranchStep{}
	case "CreateBranchStep":
		step := CreateBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "CreateAndCheckoutBranchStep":
		step := CreateAndCheckoutBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "CreatePullRequestStep":
		step := CreatePullRequestStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "CreateRemoteBranchStep":
		step := CreateRemoteBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "CreateTrackingBranchStep":
		step := CreateTrackingBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "DeleteAncestorBranchesStep":
		return DeleteAncestorBranchesStep{}
	case "DeleteLocalBranchStep":
		step := DeleteLocalBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "DeleteParentBranchStep":
		step := DeleteParentBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "DeleteRemoteBranchStep":
		step := DeleteRemoteBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "EnsureHasShippableChangesStep":
		step := EnsureHasShippableChangesStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "MergeBranchStep":
		step := MergeBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "MergeTrackingBranchStep":
		return MergeTrackingBranchStep{}
	case "NoOpStep":
		return NoOpStep{}
	case "PushBranchAfterCurrentBranchSteps":
		return PushBranchAfterCurrentBranchSteps{}
	case "PushBranchStep":
		step := PushBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "PushTagsStep":
		return PushTagsStep{}
	case "RebaseBranchStep":
		step := RebaseBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "RebaseTrackingBranchStep":
		return RebaseTrackingBranchStep{}
	case "ResetToShaStep":
		step := ResetToShaStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "RestoreOpenChangesStep":
		return RestoreOpenChangesStep{}
	case "RevertCommitStep":
		return RevertCommitStep{}
	case "SetParentBranchStep":
		step := SetParentBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "SquashMergeBranchStep":
		step := SquashMergeBranchStep{}
		json.Unmarshal(serializedStep.Data, &step)
		return step
	case "SkipCurrentBranchSteps":
		return SkipCurrentBranchSteps{}
	case "StashOpenChangesStep":
		return StashOpenChangesStep{}
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
