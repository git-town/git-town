package steps_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/src/steps"
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
	assert.NoError(t, err)
	newRunState := &steps.RunState{}
	err = json.Unmarshal(data, &newRunState)
	assert.NoError(t, err)
	assert.Equal(t, runState, newRunState)
}
