package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/shoenig/test/must"
)

func TestRunState(t *testing.T) {
	t.Parallel()

	t.Run("Marshal and Unmarshal", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			Command: "sync",
			AbortProgram: program.Program{
				&opcode.ResetCurrentBranchToSHA{
					MustHaveSHA: domain.NewSHA("222222"),
					SetToSHA:    domain.NewSHA("111111"),
					Hard:        false,
				},
			},
			RunProgram: program.Program{
				&opcode.ResetCurrentBranchToSHA{
					MustHaveSHA: domain.NewSHA("222222"),
					SetToSHA:    domain.NewSHA("111111"),
					Hard:        false,
				},
			},
			UndoProgram: program.Program{
				&opcode.ResetCurrentBranchToSHA{
					MustHaveSHA: domain.NewSHA("222222"),
					SetToSHA:    domain.NewSHA("111111"),
					Hard:        false,
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
  "IsUndo": false,
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
  "InitialActiveBranch": "initial",
  "FinalUndoProgram": [],
  "UnfinishedDetails": null,
  "UndoablePerennialCommits": []
}`[1:]
		must.EqOp(t, want, string(encoded))
		newRunState := runstate.EmptyRunState()
		err = json.Unmarshal(encoded, &newRunState)
		must.NoError(t, err)
		must.Eq(t, runState, &newRunState)
	})
}
