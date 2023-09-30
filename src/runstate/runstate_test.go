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
			AbortSteps: runstate.StepList{
				List: []steps.Step{
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        false,
					},
				},
			},
			RunSteps: runstate.StepList{
				List: []steps.Step{
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        false,
					},
				},
			},
			UndoSteps: runstate.StepList{
				List: []steps.Step{
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        false,
					},
				},
			},
			UndoablePerennialCommits: []domain.SHA{},
			InitialActiveBranch:      domain.NewLocalBranchName("initial"),
		}
		encoded, err := json.MarshalIndent(runState, "", "  ")
		assert.NoError(t, err)
		want := `
{
  "Command": "sync",
  "IsAbort": false,
  "IsUndo": false,
  "AbortStepList": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHAStep"
    }
  ],
  "RunStepList": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHAStep"
    }
  ],
  "UndoStepList": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHAStep"
    }
  ],
  "InitialActiveBranch": "initial",
  "FinalUndoStepList": [],
  "UnfinishedDetails": null,
  "UndoablePerennialCommits": []
}`[1:]
		assert.Equal(t, want, string(encoded))
		newRunState := &runstate.RunState{} //nolint:exhaustruct
		err = json.Unmarshal(encoded, &newRunState)
		assert.NoError(t, err)
		assert.Equal(t, runState, newRunState)
	})
}
