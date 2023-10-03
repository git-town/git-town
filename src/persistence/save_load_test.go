package persistence_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/shoenig/test"
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
			rootDir := domain.NewRepoRootDir(give)
			have := persistence.SanitizePath(rootDir)
			test.EqOp(t, want, have)
		}
	})

	t.Run("Save and Load", func(t *testing.T) {
		t.Parallel()
		runState := runstate.RunState{
			Command:    "command",
			IsAbort:    true,
			IsUndo:     true,
			AbortSteps: runstate.StepList{},
			RunSteps: runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.AbortRebaseStep{},
					&steps.AddToPerennialBranchesStep{Branch: domain.NewLocalBranchName("branch")},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch")},
					&steps.CommitOpenChangesStep{},
					&steps.ConnectorMergeProposalStep{
						Branch:          domain.NewLocalBranchName("branch"),
						CommitMessage:   "commit message",
						ProposalMessage: "proposal message",
						ProposalNumber:  123,
					},
					&steps.ContinueMergeStep{},
					&steps.ContinueRebaseStep{},
					&steps.CreateBranchStep{
						Branch:        domain.NewLocalBranchName("branch"),
						StartingPoint: domain.NewSHA("123456").Location(),
					},
					&steps.CreateProposalStep{Branch: domain.NewLocalBranchName("branch")},
					&steps.CreateRemoteBranchStep{
						Branch:     domain.NewLocalBranchName("branch"),
						NoPushHook: true,
						SHA:        domain.NewSHA("123456"),
					},
					&steps.CreateTrackingBranchStep{
						Branch:     domain.NewLocalBranchName("branch"),
						NoPushHook: true,
					},
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent").Location(),
						Force:  false,
					},
					&steps.DeleteRemoteBranchStep{
						Branch: domain.NewRemoteBranchName("origin/branch"),
					},
					&steps.DeleteParentBranchStep{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&steps.DeleteTrackingBranchStep{
						Branch: domain.NewRemoteBranchName("origin/branch"),
					},
					&steps.DiscardOpenChangesStep{},
					&steps.EnsureHasShippableChangesStep{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent"),
					},
					&steps.FetchUpstreamStep{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&steps.ForcePushCurrentBranchStep{
						NoPushHook: true,
					},
					&steps.MergeStep{Branch: domain.NewBranchName("branch")},
					&steps.PreserveCheckoutHistoryStep{
						InitialBranch:                     domain.NewLocalBranchName("initial-branch"),
						InitialPreviouslyCheckedOutBranch: domain.NewLocalBranchName("initial-previous-branch"),
						MainBranch:                        domain.NewLocalBranchName("main"),
					},
					&steps.PullCurrentBranchStep{},
					&steps.PushCurrentBranchStep{
						CurrentBranch: domain.NewLocalBranchName("branch"),
						NoPushHook:    true,
					},
					&steps.PushTagsStep{},
					&steps.RebaseBranchStep{Branch: domain.NewBranchName("branch")},
					&steps.RemoveFromPerennialBranchesStep{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&steps.RemoveGlobalConfigStep{
						Key: config.KeyOffline,
					},
					&steps.RemoveLocalConfigStep{
						Key: config.KeyOffline,
					},
					&steps.ResetCurrentBranchToSHAStep{
						Hard:        true,
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
					},
					&steps.RestoreOpenChangesStep{},
					&steps.RevertCommitStep{
						SHA: domain.NewSHA("123456"),
					},
					&steps.SetGlobalConfigStep{
						Key:   config.KeyOffline,
						Value: "1",
					},
					&steps.SetLocalConfigStep{
						Key:   config.KeyOffline,
						Value: "1",
					},
					&steps.SetParentStep{
						Branch:       domain.NewLocalBranchName("branch"),
						ParentBranch: domain.NewLocalBranchName("parent"),
					},
					&steps.SkipCurrentBranchSteps{},
					&steps.SquashMergeStep{
						Branch:        domain.NewLocalBranchName("branch"),
						CommitMessage: "commit message",
						Parent:        domain.NewLocalBranchName("parent"),
					},
					&steps.StashOpenChangesStep{},
					&steps.UpdateProposalTargetStep{
						ProposalNumber: 123,
						NewTarget:      domain.NewLocalBranchName("new-target"),
					},
				},
			},
			UndoSteps: runstate.StepList{},
			UnfinishedDetails: &runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: domain.NewLocalBranchName("end-branch"),
				EndTime:   time.Time{},
			},
			InitialActiveBranch:      domain.NewLocalBranchName("initial"),
			UndoablePerennialCommits: []domain.SHA{},
		}

		wantJSON := `
{
  "Command": "command",
  "IsAbort": true,
  "IsUndo": true,
  "AbortSteps": [],
  "RunSteps": [
    {
      "data": {},
      "type": "AbortMergeStep"
    },
    {
      "data": {},
      "type": "AbortRebaseStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToPerennialBranchesStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CheckoutStep"
    },
    {
      "data": {},
      "type": "CommitOpenChangesStep"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "ProposalMessage": "proposal message",
        "ProposalNumber": 123
      },
      "type": "ConnectorMergeProposalStep"
    },
    {
      "data": {},
      "type": "ContinueMergeStep"
    },
    {
      "data": {},
      "type": "ContinueRebaseStep"
    },
    {
      "data": {
        "Branch": "branch",
        "StartingPoint": "123456"
      },
      "type": "CreateBranchStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CreateProposalStep"
    },
    {
      "data": {
        "Branch": "branch",
        "NoPushHook": true,
        "SHA": "123456"
      },
      "type": "CreateRemoteBranchStep"
    },
    {
      "data": {
        "Branch": "branch",
        "NoPushHook": true
      },
      "type": "CreateTrackingBranchStep"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent",
        "Force": false
      },
      "type": "DeleteLocalBranchStep"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "DeleteRemoteBranchStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "DeleteParentBranchStep"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "DeleteTrackingBranchStep"
    },
    {
      "data": {},
      "type": "DiscardOpenChangesStep"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "EnsureHasShippableChangesStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "FetchUpstreamStep"
    },
    {
      "data": {
        "NoPushHook": true
      },
      "type": "ForcePushCurrentBranchStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "MergeStep"
    },
    {
      "data": {
        "InitialBranch": "initial-branch",
        "InitialPreviouslyCheckedOutBranch": "initial-previous-branch",
        "MainBranch": "main"
      },
      "type": "PreserveCheckoutHistoryStep"
    },
    {
      "data": {},
      "type": "PullCurrentBranchStep"
    },
    {
      "data": {
        "CurrentBranch": "branch",
        "NoPushHook": true
      },
      "type": "PushCurrentBranchStep"
    },
    {
      "data": {},
      "type": "PushTagsStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RebaseBranchStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromPerennialBranchesStep"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "RemoveGlobalConfigStep"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "RemoveLocalConfigStep"
    },
    {
      "data": {
        "Hard": true,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHAStep"
    },
    {
      "data": {},
      "type": "RestoreOpenChangesStep"
    },
    {
      "data": {
        "SHA": "123456"
      },
      "type": "RevertCommitStep"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "SetGlobalConfigStep"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "SetLocalConfigStep"
    },
    {
      "data": {
        "Branch": "branch",
        "ParentBranch": "parent"
      },
      "type": "SetParentStep"
    },
    {
      "data": {},
      "type": "SkipCurrentBranchSteps"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Parent": "parent"
      },
      "type": "SquashMergeStep"
    },
    {
      "data": {},
      "type": "StashOpenChangesStep"
    },
    {
      "data": {
        "ProposalNumber": 123,
        "NewTarget": "new-target"
      },
      "type": "UpdateProposalTargetStep"
    }
  ],
  "UndoSteps": [],
  "InitialActiveBranch": "initial",
  "FinalUndoSteps": [],
  "UnfinishedDetails": {
    "CanSkip": true,
    "EndBranch": "end-branch",
    "EndTime": "0001-01-01T00:00:00Z"
  },
  "UndoablePerennialCommits": []
}`[1:]

		repoRoot := domain.NewRepoRootDir("/path/to/git-town-unit-tests")
		err := persistence.Save(&runState, repoRoot)
		test.NoError(t, err)
		filepath, err := persistence.FilePath(repoRoot)
		test.NoError(t, err)
		content, err := os.ReadFile(filepath)
		test.NoError(t, err)
		test.EqOp(t, wantJSON, string(content))
		var newState runstate.RunState
		err = json.Unmarshal(content, &newState)
		test.NoError(t, err)
		test.Eq(t, runState, newState)
	})
}
