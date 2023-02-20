package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestRunState(t *testing.T) {
	t.Parallel()
	t.Run(".Marshal()", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			AbortStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
			},
			Command: "sync",
			RunStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
			},
			UndoStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
			},
		}
		data, err := json.Marshal(runState)
		assert.NoError(t, err)
		newRunState := &runstate.RunState{}
		err = json.Unmarshal(data, &newRunState)
		assert.NoError(t, err)
		assert.Equal(t, runState, newRunState)
	})
}
