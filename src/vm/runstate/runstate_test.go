package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
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
			EndBranchesSnapshot: gitdomain.BranchesSnapshot{
				Active: Some(gitdomain.NewLocalBranchName("branch-1")),
				Branches: gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
						SyncStatus: gitdomain.SyncStatusNotInSync,
					},
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-2")),
						LocalSHA:   Some(gitdomain.NewSHA("333333")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					},
				},
			},
			EndConfigSnapshot:        undoconfig.EmptyConfigSnapshot(),
			EndStashSize:             1,
			BeginBranchesSnapshot:    gitdomain.EmptyBranchesSnapshot(),
			BeginConfigSnapshot:      undoconfig.EmptyConfigSnapshot(),
			BeginStashSize:           0,
			UndoablePerennialCommits: []gitdomain.SHA{},
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
  "BeginBranchesSnapshot": {
    "Active": null,
    "Branches": []
  },
  "BeginConfigSnapshot": {
    "Global": {},
    "Local": {}
  },
  "BeginStashSize": 0,
  "Command": "sync",
  "DryRun": true,
  "EndBranchesSnapshot": {
    "Active": "branch-1",
    "Branches": [
      {
        "LocalName": "branch-1",
        "LocalSHA": "111111",
        "RemoteName": "origin/branch-1",
        "RemoteSHA": "222222",
        "SyncStatus": "not in sync"
      },
      {
        "LocalName": "branch-2",
        "LocalSHA": "333333",
        "RemoteName": null,
        "RemoteSHA": null,
        "SyncStatus": "local only"
      }
    ]
  },
  "EndConfigSnapshot": {
    "Global": {},
    "Local": {}
  },
  "EndStashSize": 1,
  "FinalUndoProgram": [],
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
