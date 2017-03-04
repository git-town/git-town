package steps

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "reflect"
)


func saveState(runState RunState) {
  serializedRunState := SerializedRunState{
    AbortStep: serializeStep(runState.AbortStep),
    RunSteps: serializeSteps(runState.RunStepList.List),
    UndoSteps: serializeSteps(runState.UndoStepList.List),
  }
  content, err := json.Marshal(serializedRunState)
  if err != nil {
    log.Fatal(err)
  }
  filename := getRunResultFilename(runState.Command)
  err = ioutil.WriteFile(filename, content, 0644)
  if err != nil {
    log.Fatal(err)
  }
}


func serializeStep(step Step) SerializedStep {
  data, err := json.Marshal(step)
  if err != nil {
    log.Fatal(err)
  }
  return SerializedStep{
    Data: data,
    Type: getType(step),
  }
}


func serializeSteps(steps []Step) []SerializedStep {
  var output []SerializedStep
  for _, step := range(steps) {
    output = append(output, serializeStep(step))
  }
  return output
}


func getType(myvar interface{}) string {
  if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
    return "*" + t.Elem().Name()
  } else {
    return t.Name()
  }
}
