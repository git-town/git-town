package statefile_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/vm/opcode"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/statefile"
	"github.com/shoenig/test/must"
)

func TestLoadSave(t *testing.T) {
	t.Parallel()

	t.Run("SanitizePath", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"/home/user/development/git-town":        "home-user-development-git-town",
			"c:\\Users\\user\\development\\git-town": "c-users-user-development-git-town",
		}
		for give, want := range tests {
			rootDir := gitdomain.NewRepoRootDir(give)
			have := statefile.SanitizePath(rootDir)
			must.EqOp(t, want, have)
		}
	})

	t.Run("Save and Load", func(t *testing.T) {
		t.Parallel()
		runState := runstate.RunState{
			Command:                "command",
			IsUndo:                 true,
			AbortProgram:           program.Program{},
			AfterBranchesSnapshot:  gitdomain.EmptyBranchesSnapshot(),
			AfterConfigSnapshot:    undoconfig.EmptyConfigSnapshot(),
			AfterStashSize:         1,
			BeforeBranchesSnapshot: gitdomain.EmptyBranchesSnapshot(),
			BeforeConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
			BeforeStashSize:        0,
			DryRun:                 true,
			RunProgram: program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
				&opcode.AddToPerennialBranches{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcode.ChangeParent{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcode.CommitOpenChanges{},
				&opcode.ConnectorMergeProposal{
					Branch:          gitdomain.NewLocalBranchName("branch"),
					CommitMessage:   "commit message",
					ProposalMessage: "proposal message",
					ProposalNumber:  123,
				},
				&opcode.ContinueMerge{},
				&opcode.ContinueRebase{},
				&opcode.CreateBranch{
					Branch:        gitdomain.NewLocalBranchName("branch"),
					StartingPoint: gitdomain.NewSHA("123456").Location(),
				},
				&opcode.CreateProposal{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcode.CreateRemoteBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
					SHA:    gitdomain.NewSHA("123456"),
				},
				&opcode.CreateTrackingBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcode.DeleteLocalBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Force:  false,
				},
				&opcode.DeleteParentBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcode.DeleteTrackingBranch{
					Branch: gitdomain.NewRemoteBranchName("origin/branch"),
				},
				&opcode.DiscardOpenChanges{},
				&opcode.EndOfBranchProgram{},
				&opcode.EnsureHasShippableChanges{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcode.FetchUpstream{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcode.ForcePushCurrentBranch{},
				&opcode.Merge{Branch: gitdomain.NewBranchName("branch")},
				&opcode.MergeParent{
					CurrentBranch:               gitdomain.NewLocalBranchName("branch"),
					ParentActiveInOtherWorktree: true,
				},
				&opcode.PreserveCheckoutHistory{
					PreviousBranchCandidates: gitdomain.NewLocalBranchNames("previous"),
				},
				&opcode.PullCurrentBranch{},
				&opcode.PushCurrentBranch{
					CurrentBranch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcode.PushTags{},
				&opcode.RebaseBranch{Branch: gitdomain.NewBranchName("branch")},
				&opcode.RebaseParent{
					CurrentBranch:               gitdomain.NewLocalBranchName("branch"),
					ParentActiveInOtherWorktree: true,
				},
				&opcode.RemoveFromPerennialBranches{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcode.RemoveGlobalConfig{
					Key: gitconfig.KeyOffline,
				},
				&opcode.RemoveLocalConfig{
					Key: gitconfig.KeyOffline,
				},
				&opcode.ResetCurrentBranchToSHA{
					Hard:        true,
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
				},
				&opcode.RestoreOpenChanges{},
				&opcode.RevertCommit{
					SHA: gitdomain.NewSHA("123456"),
				},
				&opcode.SetGlobalConfig{
					Key:   gitconfig.KeyOffline,
					Value: "1",
				},
				&opcode.SetLocalConfig{
					Key:   gitconfig.KeyOffline,
					Value: "1",
				},
				&opcode.SetParent{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcode.SetParentIfBranchExists{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcode.SkipCurrentBranch{},
				&opcode.SquashMerge{
					Branch:        gitdomain.NewLocalBranchName("branch"),
					CommitMessage: "commit message",
					Parent:        gitdomain.NewLocalBranchName("parent"),
				},
				&opcode.StashOpenChanges{},
				&opcode.UpdateProposalTarget{
					ProposalNumber: 123,
					NewTarget:      gitdomain.NewLocalBranchName("new-target"),
				},
			},
			UndoProgram: program.Program{},
			UnfinishedDetails: &runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: gitdomain.NewLocalBranchName("end-branch"),
				EndTime:   time.Time{},
			},
			UndoablePerennialCommits: []gitdomain.SHA{},
		}

		wantJSON := `
{
  "AbortProgram": [],
  "AfterBranchesSnapshot": {
    "Active": "",
    "Branches": []
  },
  "AfterConfigSnapshot": {
    "Global": {},
    "Local": {}
  },
  "AfterStashSize": 1,
  "BeforeBranchesSnapshot": {
    "Active": "",
    "Branches": []
  },
  "BeforeConfigSnapshot": {
    "Global": {},
    "Local": {}
  },
  "BeforeStashSize": 0,
  "Command": "command",
  "DryRun": true,
  "FinalUndoProgram": [],
  "IsUndo": true,
  "RunProgram": [
    {
      "data": {},
      "type": "AbortMerge"
    },
    {
      "data": {},
      "type": "AbortRebase"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToPerennialBranches"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "ChangeParent"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "Checkout"
    },
    {
      "data": {},
      "type": "CommitOpenChanges"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "ProposalMessage": "proposal message",
        "ProposalNumber": 123
      },
      "type": "ConnectorMergeProposal"
    },
    {
      "data": {},
      "type": "ContinueMerge"
    },
    {
      "data": {},
      "type": "ContinueRebase"
    },
    {
      "data": {
        "Branch": "branch",
        "StartingPoint": "123456"
      },
      "type": "CreateBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CreateProposal"
    },
    {
      "data": {
        "Branch": "branch",
        "SHA": "123456"
      },
      "type": "CreateRemoteBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CreateTrackingBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "Force": false
      },
      "type": "DeleteLocalBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "DeleteParentBranch"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "DeleteTrackingBranch"
    },
    {
      "data": {},
      "type": "DiscardOpenChanges"
    },
    {
      "data": {},
      "type": "EndOfBranchProgram"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "EnsureHasShippableChanges"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "FetchUpstream"
    },
    {
      "data": {},
      "type": "ForcePushCurrentBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "Merge"
    },
    {
      "data": {
        "CurrentBranch": "branch",
        "ParentActiveInOtherWorktree": true
      },
      "type": "MergeParent"
    },
    {
      "data": {
        "PreviousBranchCandidates": [
          "previous"
        ]
      },
      "type": "PreserveCheckoutHistory"
    },
    {
      "data": {},
      "type": "PullCurrentBranch"
    },
    {
      "data": {
        "CurrentBranch": "branch"
      },
      "type": "PushCurrentBranch"
    },
    {
      "data": {},
      "type": "PushTags"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RebaseBranch"
    },
    {
      "data": {
        "CurrentBranch": "branch",
        "ParentActiveInOtherWorktree": true
      },
      "type": "RebaseParent"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromPerennialBranches"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "RemoveGlobalConfig"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "RemoveLocalConfig"
    },
    {
      "data": {
        "Hard": true,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    },
    {
      "data": {},
      "type": "RestoreOpenChanges"
    },
    {
      "data": {
        "SHA": "123456"
      },
      "type": "RevertCommit"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "SetGlobalConfig"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "SetLocalConfig"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "SetParent"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "SetParentIfBranchExists"
    },
    {
      "data": {},
      "type": "SkipCurrentBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Parent": "parent"
      },
      "type": "SquashMerge"
    },
    {
      "data": {},
      "type": "StashOpenChanges"
    },
    {
      "data": {
        "NewTarget": "new-target",
        "ProposalNumber": 123
      },
      "type": "UpdateProposalTarget"
    }
  ],
  "UndoProgram": [],
  "UndoablePerennialCommits": [],
  "UnfinishedDetails": {
    "CanSkip": true,
    "EndBranch": "end-branch",
    "EndTime": "0001-01-01T00:00:00Z"
  }
}`[1:]

		repoRoot := gitdomain.NewRepoRootDir("/path/to/git-town-unit-tests")
		err := statefile.Save(&runState, repoRoot)
		must.NoError(t, err)
		filepath, err := statefile.FilePath(repoRoot)
		must.NoError(t, err)
		content, err := os.ReadFile(filepath)
		must.NoError(t, err)
		must.EqOp(t, wantJSON, string(content))
		var newState runstate.RunState
		err = json.Unmarshal(content, &newState)
		must.NoError(t, err)
		// NOTE: comparing runState and newState directly leads to incorrect test failures
		// solely due to different pointer addresses, even when using reflect.DeepEqual.
		// Comparing the serialization seems to work better here.
		runStateText := fmt.Sprintf("%+v", runState)
		newStateText := fmt.Sprintf("%+v", newState)
		must.EqOp(t, runStateText, newStateText)
	})
}
