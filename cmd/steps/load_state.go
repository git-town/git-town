package steps

import(
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)


func loadState(command string) RunState {
  filename := getRunResultFilename(command)
  content, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }
  var serializedRunState SerializedRunState
  err = json.Unmarshal(content, &serializedRunState)
  if err != nil {
    log.Fatal(err)
  }
  return RunState{
    AbortStep: deserializeStep(serializedRunState.AbortStep),
    RunStepList: deserializeStepList(serializedRunState.RunSteps),
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
  case "CreateAndCheckoutBranchStep":
    step := CreateAndCheckoutBranchStep{}
    json.Unmarshal(serializedStep.Data, &step)
    return step
  case "CreateTrackingBranchStep":
    step := CreateTrackingBranchStep{}
    json.Unmarshal(serializedStep.Data, &step)
    return step
  case "MergeBranchStep":
    step := MergeBranchStep{}
    json.Unmarshal(serializedStep.Data, &step)
    return step
  case "MergeTrackingBranchStep":
    return MergeTrackingBranchStep{}
  case "PushBranchStep":
    step := PushBranchStep{}
    json.Unmarshal(serializedStep.Data, &step)
    return step
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
  case "StashOpenChangesStep":
    return StashOpenChangesStep{}
  }
  log.Fatal(fmt.Sprintf("Cannot deserialize steps: %s %s", serializedStep.Type, serializedStep.Data))
  return nil
}


func deserializeStepList(serializedSteps []SerializedStep) (stepList StepList) {
  for i := 0; i < len(serializedSteps); i++ {
    stepList.Append(deserializeStep(serializedSteps[i]))
  }
  return stepList
}
