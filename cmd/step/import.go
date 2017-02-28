package step

import(
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

func Import(commandName string) RunResult {
  filename := getRunResultFilename(commandName)
  content, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }
  var runResultData SerializedRunResult
  err = json.Unmarshal(content, &runResultData)
  if err != nil {
    log.Fatal(err)
  }
  return RunResult{
    AbortStep: importStep(runResultData.AbortStep),
    ContinueSteps: importSteps(runResultData.ContinueSteps),
    UndoSteps: importSteps(runResultData.UndoSteps),
  }
}

func importStep(serializedStep SerializedStep) Step {
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

func importSteps(serializedSteps []SerializedStep) []Step {
  var output []Step
  for i := 0; i < len(serializedSteps); i++ {
    output = append(output, importStep(serializedSteps[i]))
  }
  return output
}
