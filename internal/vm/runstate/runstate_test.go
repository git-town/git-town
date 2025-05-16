package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/program"
	"github.com/git-town/git-town/v20/internal/vm/runstate"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestRunState(t *testing.T) {
	t.Parallel()

	t.Run("Marshal and Unmarshal", func(t *testing.T) {
		t.Parallel()
		runState := &runstate.RunState{
			BranchInfosLastRun: Some(gitdomain.BranchInfos{
				{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusAhead,
				},
			}),
			Command: "sync",
			DryRun:  true,
			AbortProgram: program.Program{
				&opcodes.BranchCurrentResetToSHAIfNeeded{
					MustHaveSHA: "222222",
					SetToSHA:    "111111",
				},
			},
			RunProgram: program.Program{
				&opcodes.BranchCurrentResetToSHAIfNeeded{
					MustHaveSHA: "222222",
					SetToSHA:    "111111",
				},
			},
			EndBranchesSnapshot: Some(gitdomain.BranchesSnapshot{
				Active: gitdomain.NewLocalBranchNameOption("branch-1"),
				Branches: gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
						SyncStatus: gitdomain.SyncStatusNotInSync,
					},
					gitdomain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchNameOption("branch-2"),
						LocalSHA:   Some(gitdomain.NewSHA("333333")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					},
				},
			}),
			EndConfigSnapshot:        None[undoconfig.ConfigSnapshot](),
			EndStashSize:             Some(gitdomain.StashSize(1)),
			BeginBranchesSnapshot:    gitdomain.EmptyBranchesSnapshot(),
			BeginConfigSnapshot:      undoconfig.EmptyConfigSnapshot(),
			BeginStashSize:           0,
			UndoablePerennialCommits: []gitdomain.SHA{},
			TouchedBranches:          []gitdomain.BranchName{"branch-1", "branch-2"},
		}
		encoded, err := json.MarshalIndent(runState, "", "  ")
		must.NoError(t, err)
		want := `
{
  "AbortProgram": [
    {
      "data": {
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "BranchCurrentResetToSHAIfNeeded"
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
  "BranchInfosLastRun": [
    {
      "LocalName": "branch",
      "LocalSHA": "111111",
      "RemoteName": "origin/branch",
      "RemoteSHA": "222222",
      "SyncStatus": "ahead"
    }
  ],
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
  "EndConfigSnapshot": null,
  "EndStashSize": 1,
  "FinalUndoProgram": [],
  "RunProgram": [
    {
      "data": {
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "BranchCurrentResetToSHAIfNeeded"
    }
  ],
  "TouchedBranches": [
    "branch-1",
    "branch-2"
  ],
  "UndoAPIProgram": [],
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
