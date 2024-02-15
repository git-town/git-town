package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/shoenig/test/must"
)

func TestRunState(t *testing.T) {
	t.Parallel()

	t.Run("Marshal and Unmarshal", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			Command: "sync",
			DryRun:  true,
			AbortProgram: program.Program{
				&opcodes.ResetCurrentBranchToSHA{
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
					Hard:        false,
				},
			},
			RunProgram: program.Program{
				&opcodes.ResetCurrentBranchToSHA{
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
					Hard:        false,
				},
			},
			UndoProgram: program.Program{
				&opcodes.ResetCurrentBranchToSHA{
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
					Hard:        false,
				},
			},
			UndoablePerennialCommits: []gitdomain.SHA{},
			InitialActiveBranch:      gitdomain.NewLocalBranchName("initial"),
		}
		encoded, err := json.MarshalIndent(runState, "", "  ")
		must.NoError(t, err)
		want := `
{
  "AbortProgram": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    }
  ],
  "Command": "sync",
  "DryRun": true,
  "FinalUndoProgram": [],
  "InitialActiveBranch": "initial",
  "IsUndo": false,
  "RunProgram": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    }
  ],
  "UndoProgram": [
    {
      "data": {
        "Hard": false,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    }
  ],
  "UndoablePerennialCommits": [],
  "UnfinishedDetails": null
}`[1:]
		must.EqOp(t, want, string(encoded))
		newRunState := runstate.EmptyRunState()
		err = json.Unmarshal(encoded, &newRunState)
		must.NoError(t, err)
		must.Eq(t, runState, &newRunState)
	})
}
