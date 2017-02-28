package step

import(
  "encoding/json"
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
    AbortSteps: importSteps(runResultData.AbortSteps),
    ContinueSteps: importSteps(runResultData.ContinueSteps),
    Success: runResultData.Success,
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
  case "RestoreOpenChangesStep":
    return RestoreOpenChangesStep{}
  }
  return nil
}

func importSteps(serializedSteps []SerializedStep) []Step {
  var output []Step
  for i := 0; i < len(serializedSteps); i++ {
    output = append(output, importStep(serializedSteps[i]))
  }
  return output
}
