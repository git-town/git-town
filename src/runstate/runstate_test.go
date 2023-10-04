package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/shoenig/test/must"
)

func TestRunState(t *testing.T) {
	t.Parallel()

	t.Run("Marshal and Unmarshal", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			Command: "sync",
			AbortSteps: runstate.StepList{
				List: []step.Step{
					&step.ResetCurrentBranchToSHA{
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        false,
					},
				},
			},
			RunSteps: runstate.StepList{
				List: []step.Step{
					&step.ResetCurrentBranchToSHA{
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        false,
					},
				},
			},
			UndoSteps: runstate.StepList{
				List: []step.Step{
					&step.ResetCurrentBranchToSHA{
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
		must.NoError(t, err)
		want := `
{
  "Command": "sync",
  "IsAbort": false,
  "IsUndo": false,
  "AbortSteps": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    }
  ],
  "RunSteps": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    }
  ],
  "UndoSteps": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    }
  ],
  "InitialActiveBranch": "initial",
  "FinalUndoSteps": [],
  "UnfinishedDetails": null,
  "UndoablePerennialCommits": []
}`[1:]
		must.EqOp(t, want, string(encoded))
		newRunState := &runstate.RunState{} //nolint:exhaustruct
		err = json.Unmarshal(encoded, &newRunState)
		must.NoError(t, err)
		must.Eq(t, runState, newRunState)
	})
}
