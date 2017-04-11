package steps

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
)

func saveState(runState *RunState) {
	serializedRunState := SerializedRunState{
		AbortStep: serializeStep(runState.AbortStep),
		RunSteps:  serializeSteps(runState.RunStepList.List),
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
		Type: getTypeName(step),
	}
}

func serializeSteps(steps []Step) (result []SerializedStep) {
	for _, step := range steps {
		result = append(result, serializeStep(step))
	}
	return
}

func getTypeName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}
