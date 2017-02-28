package step

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "reflect"
)

func export(commandName string, abortSteps, continueSteps []Step) {
  runResultData := SerializedRunResult{
    AbortSteps: exportSteps(abortSteps),
    ContinueSteps: exportSteps(continueSteps),
    Success: false,
  }
  content, err := json.Marshal(runResultData)
  if err != nil {
    log.Fatal(err)
  }
  filename := getRunResultFilename(commandName)
  err = ioutil.WriteFile(filename, content, 0644)
  if err != nil {
    log.Fatal(err)
  }
}

func getType(myvar interface{}) string {
  if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
    return "*" + t.Elem().Name()
  } else {
    return t.Name()
  }
}

func exportStep(step Step) SerializedStep {
  data, err := json.Marshal(step)
  if err != nil {
    log.Fatal(err)
  }
  return SerializedStep{
    Data: data,
    Type: getType(step),
  }
}

func exportSteps(steps []Step) []SerializedStep {
  var output []SerializedStep
  for i := 0; i < len(steps); i++ {
    step := steps[i]
    output = append(output, exportStep(step))
  }
  return output
}
