package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcode"
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
				&opcode.ResetCurrentBranchToSHA{
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
					Hard:        false,
				},
			},
			RunProgram: program.Program{
				&opcode.ResetCurrentBranchToSHA{
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
					Hard:        false,
				},
			},
			UndoProgram: program.Program{
				&opcode.ResetCurrentBranchToSHA{
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
					Hard:        false,
				},
			},
			AfterBranchesSnapshot: gitdomain.BranchesSnapshot{
				Active: "branch-1",
				Branches: gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  "branch-1",
						LocalSHA:   "111111",
						RemoteName: "origin/branch-1",
						RemoteSHA:  "222222",
						SyncStatus: gitdomain.SyncStatusNotInSync,
					},
					gitdomain.BranchInfo{
						LocalName:  "branch-2",
						LocalSHA:   "333333",
						RemoteName: gitdomain.EmptyRemoteBranchName(),
						RemoteSHA:  gitdomain.EmptySHA(),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					},
				},
			},
			AfterGlobalConfigSnapshot:  configdomain.EmptyPartialConfig(),
			AfterLocalConfigSnapshot:   configdomain.EmptyPartialConfig(),
			BeforeBranchesSnapshot:     gitdomain.EmptyBranchesSnapshot(),
			BeforeGlobalConfigSnapshot: configdomain.EmptyPartialConfig(),
			BeforeLocalConfigSnapshot:  configdomain.EmptyPartialConfig(),
			UndoablePerennialCommits:   []gitdomain.SHA{},
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
  "AfterBranchesSnapshot": {
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
        "RemoteName": "",
        "RemoteSHA": "",
        "SyncStatus": "local only"
      }
    ]
  },
  "AfterGlobalConfigSnapshot": {
    "Aliases": {},
    "GitHubToken": null,
    "GitLabToken": null,
    "GitUserEmail": null,
    "GitUserName": null,
    "GiteaToken": null,
    "HostingOriginHostname": null,
    "HostingPlatform": null,
    "Lineage": null,
    "MainBranch": null,
    "Offline": null,
    "PerennialBranches": null,
    "PushHook": null,
    "PushNewBranches": null,
    "ShipDeleteTrackingBranch": null,
    "SyncBeforeShip": null,
    "SyncFeatureStrategy": null,
    "SyncPerennialStrategy": null,
    "SyncUpstream": null
  },
  "AfterLocalConfigSnapshot": {
    "Aliases": {},
    "GitHubToken": null,
    "GitLabToken": null,
    "GitUserEmail": null,
    "GitUserName": null,
    "GiteaToken": null,
    "HostingOriginHostname": null,
    "HostingPlatform": null,
    "Lineage": null,
    "MainBranch": null,
    "Offline": null,
    "PerennialBranches": null,
    "PushHook": null,
    "PushNewBranches": null,
    "ShipDeleteTrackingBranch": null,
    "SyncBeforeShip": null,
    "SyncFeatureStrategy": null,
    "SyncPerennialStrategy": null,
    "SyncUpstream": null
  },
  "BeforeBranchesSnapshot": {
    "Active": "",
    "Branches": []
  },
  "BeforeGlobalConfigSnapshot": {
    "Aliases": {},
    "GitHubToken": null,
    "GitLabToken": null,
    "GitUserEmail": null,
    "GitUserName": null,
    "GiteaToken": null,
    "HostingOriginHostname": null,
    "HostingPlatform": null,
    "Lineage": null,
    "MainBranch": null,
    "Offline": null,
    "PerennialBranches": null,
    "PushHook": null,
    "PushNewBranches": null,
    "ShipDeleteTrackingBranch": null,
    "SyncBeforeShip": null,
    "SyncFeatureStrategy": null,
    "SyncPerennialStrategy": null,
    "SyncUpstream": null
  },
  "BeforeLocalConfigSnapshot": {
    "Aliases": {},
    "GitHubToken": null,
    "GitLabToken": null,
    "GitUserEmail": null,
    "GitUserName": null,
    "GiteaToken": null,
    "HostingOriginHostname": null,
    "HostingPlatform": null,
    "Lineage": null,
    "MainBranch": null,
    "Offline": null,
    "PerennialBranches": null,
    "PushHook": null,
    "PushNewBranches": null,
    "ShipDeleteTrackingBranch": null,
    "SyncBeforeShip": null,
    "SyncFeatureStrategy": null,
    "SyncPerennialStrategy": null,
    "SyncUpstream": null
  },
  "Command": "sync",
  "DryRun": true,
  "FinalUndoProgram": [],
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
