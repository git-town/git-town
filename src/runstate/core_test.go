package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestRunState(t *testing.T) {
	t.Parallel()
	t.Run(".Marshal()", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			AbortStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetToShaStep{Sha: git.NewSHA("abcdef")}},
			},
			Command: "sync",
			RunStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetToShaStep{Sha: git.NewSHA("abcdef")}},
			},
			UndoStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetToShaStep{Sha: git.NewSHA("abcdef")}},
			},
		}
		data, err := json.Marshal(runState)
		assert.NoError(t, err)
		newRunState := &runstate.RunState{} //nolint:exhaustruct
		err = json.Unmarshal(data, &newRunState)
		assert.NoError(t, err)
		assert.Equal(t, runState, newRunState)
	})
}
