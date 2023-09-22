package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestRunState(t *testing.T) {
	t.Parallel()

	t.Run("Marshal and Unmarshal", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			Command: "sync",
			AbortStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetCurrentBranchToSHAStep{SHA: domain.NewSHA("abcdef"), Hard: false}},
			},
			RunStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetCurrentBranchToSHAStep{SHA: domain.NewSHA("abcdef"), Hard: false}},
			},
			UndoStepList: runstate.StepList{
				List: []steps.Step{&steps.ResetCurrentBranchToSHAStep{SHA: domain.NewSHA("abcdef"), Hard: false}},
			},
		}
		encoded, err := json.MarshalIndent(runState, "", "  ")
		assert.NoError(t, err)
		want := `
{
  "AbortStepList": [
    {
      "data": {
        "Hard": false,
        "SHA": "abcdef"
      },
      "type": "ResetCurrentBranchToSHAStep"
    }
  ],
  "Command": "sync",
  "IsAbort": false,
  "IsUndo": false,
  "RunStepList": [
    {
      "data": {
        "Hard": false,
        "SHA": "abcdef"
      },
      "type": "ResetCurrentBranchToSHAStep"
    }
  ],
  "UndoStepList": [
    {
      "data": {
        "Hard": false,
        "SHA": "abcdef"
      },
      "type": "ResetCurrentBranchToSHAStep"
    }
  ],
  "UnfinishedDetails": null
}`[1:]
		assert.Equal(t, want, string(encoded))
		newRunState := &runstate.RunState{} //nolint:exhaustruct
		err = json.Unmarshal(encoded, &newRunState)
		assert.NoError(t, err)
		assert.Equal(t, runState, newRunState)
	})
}
