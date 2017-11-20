package steps

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/Originate/exit"
)

func saveState(runState *RunState) {
	serializedRunState := SerializedRunState{
		AbortStep: serializeStep(runState.AbortStep),
		RunSteps:  serializeSteps(runState.RunStepList.List),
		UndoSteps: serializeSteps(runState.UndoStepList.List),
	}
	content, err := json.Marshal(serializedRunState)
	exit.If(err)
	filename := getRunResultFilename(runState.Command)
	err = ioutil.WriteFile(filename, content, 0644)
	exit.If(err)
}

func serializeStep(step Step) SerializedStep {
	data, err := json.Marshal(step)
	exit.If(err)
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
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}
