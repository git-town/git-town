package steps_test

import (
	"encoding/json"
	"testing"

	"github.com/Originate/git-town/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestRunState_Marshal(t *testing.T) {
	runState := &steps.RunState{
		AbortStepList: steps.StepList{
			List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
		},
		Command: "sync",
		RunStepList: steps.StepList{
			List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
		},
		UndoStepList: steps.StepList{
			List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
		},
	}
	data, err := json.Marshal(runState)
	assert.Nil(t, err)
	newRunState := &steps.RunState{}
	err = json.Unmarshal(data, &newRunState)
	assert.Nil(t, err)
	assert.Equal(t, runState, newRunState)
}
