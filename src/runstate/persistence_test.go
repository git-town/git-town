package runstate_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestSanitizePath(t *testing.T) {
	t.Parallel()
	t.Run("SanitizePath", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"/home/user/development/git-town":        "home-user-development-git-town",
			"c:\\Users\\user\\development\\git-town": "c-users-user-development-git-town",
		}
		for give, want := range tests {
			have := runstate.SanitizePath(give)
			assert.Equal(t, want, have)
		}
	})
	t.Run("Save and Load", func(t *testing.T) {
		t.Parallel()
		runState := runstate.RunState{
			AbortStepList: runstate.StepList{},
			Command:       "command",
			IsAbort:       true,
			RunStepList: runstate.StepList{
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
						Sha:        domain.NewSHA("123456"),
					},
					&steps.CreateTrackingBranchStep{
						Branch:     domain.NewLocalBranchName("branch"),
						NoPushHook: true,
					},
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent").Location(),
					},
					&steps.DeleteOriginBranchStep{
						Branch:     domain.NewLocalBranchName("branch"),
						IsTracking: true,
						NoPushHook: true,
					},
					&steps.DeleteParentBranchStep{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent"),
					},
					&steps.DiscardOpenChangesStep{},
					&steps.EnsureHasShippableChangesStep{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent"),
					},
					&steps.FetchUpstreamStep{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&steps.MergeStep{Branch: domain.NewBranchName("branch")},
					&steps.PreserveCheckoutHistoryStep{
						InitialBranch:                     domain.NewLocalBranchName("initial-branch"),
						InitialPreviouslyCheckedOutBranch: domain.NewLocalBranchName("initial-previous-branch"),
						MainBranch:                        domain.NewLocalBranchName("main"),
					},
					&steps.PullBranchStep{Branch: "branch"},
					&steps.PushBranchAfterCurrentBranchSteps{},
					&steps.PushBranchStep{
						Branch:         domain.NewLocalBranchName("branch"),
						TrackingBranch: domain.NewRemoteBranchName("origin/branch"),
						ForceWithLease: true,
						NoPushHook:     true,
						Undoable:       true,
					},
					&steps.PushTagsStep{},
					&steps.RebaseBranchStep{Branch: domain.NewBranchName("branch")},
					&steps.RemoveFromPerennialBranchesStep{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&steps.ResetToShaStep{
						Hard: true,
						Sha:  domain.NewSHA("123456"),
					},
					&steps.RestoreOpenChangesStep{},
					&steps.RevertCommitStep{
						Sha: domain.NewSHA("123456"),
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
						ExistingTarget: domain.NewLocalBranchName("existing-target"),
					},
				},
			},
			UndoStepList: runstate.StepList{},
			UnfinishedDetails: &runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: domain.NewLocalBranchName("end-branch"),
				EndTime:   time.Time{},
			},
		}

		wantJSON := `
{
  "AbortStepList": [],
  "Command": "command",
  "IsAbort": true,
  "RunStepList": [
    {
      "data": {},
      "type": "*AbortMergeStep"
    },
    {
      "data": {},
      "type": "*AbortRebaseStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*AddToPerennialBranchesStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*CheckoutStep"
    },
    {
      "data": {},
      "type": "*CommitOpenChangesStep"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "ProposalMessage": "proposal message",
        "ProposalNumber": 123
      },
      "type": "*ConnectorMergeProposalStep"
    },
    {
      "data": {},
      "type": "*ContinueMergeStep"
    },
    {
      "data": {},
      "type": "*ContinueRebaseStep"
    },
    {
      "data": {
        "Branch": "branch",
        "StartingPoint": "123456"
      },
      "type": "*CreateBranchStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*CreateProposalStep"
    },
    {
      "data": {
        "Branch": "branch",
        "NoPushHook": true,
        "Sha": "123456"
      },
      "type": "*CreateRemoteBranchStep"
    },
    {
      "data": {
        "Branch": "branch",
        "NoPushHook": true
      },
      "type": "*CreateTrackingBranchStep"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent",
        "Force": false
      },
      "type": "*DeleteLocalBranchStep"
    },
    {
      "data": {
        "Branch": "branch",
        "IsTracking": true,
        "NoPushHook": true
      },
      "type": "*DeleteOriginBranchStep"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "*DeleteParentBranchStep"
    },
    {
      "data": {},
      "type": "*DiscardOpenChangesStep"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "*EnsureHasShippableChangesStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*FetchUpstreamStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*MergeStep"
    },
    {
      "data": {
        "InitialBranch": "initial-branch",
        "InitialPreviouslyCheckedOutBranch": "initial-previous-branch",
        "MainBranch": "main"
      },
      "type": "*PreserveCheckoutHistoryStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*PullBranchStep"
    },
    {
      "data": {},
      "type": "*PushBranchAfterCurrentBranchSteps"
    },
    {
      "data": {
        "Branch": "branch",
        "TrackingBranch": "origin/branch",
        "ForceWithLease": true,
        "NoPushHook": true,
        "Undoable": true
      },
      "type": "*PushBranchStep"
    },
    {
      "data": {},
      "type": "*PushTagsStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*RebaseBranchStep"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "*RemoveFromPerennialBranchesStep"
    },
    {
      "data": {
        "Hard": true,
        "Sha": "123456"
      },
      "type": "*ResetToShaStep"
    },
    {
      "data": {},
      "type": "*RestoreOpenChangesStep"
    },
    {
      "data": {
        "Sha": "123456"
      },
      "type": "*RevertCommitStep"
    },
    {
      "data": {
        "Branch": "branch",
        "ParentBranch": "parent"
      },
      "type": "*SetParentStep"
    },
    {
      "data": {},
      "type": "*SkipCurrentBranchSteps"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Parent": "parent"
      },
      "type": "*SquashMergeStep"
    },
    {
      "data": {},
      "type": "*StashOpenChangesStep"
    },
    {
      "data": {
        "ProposalNumber": 123,
        "NewTarget": "new-target",
        "ExistingTarget": "existing-target"
      },
      "type": "*UpdateProposalTargetStep"
    }
  ],
  "UndoStepList": [],
  "UnfinishedDetails": {
    "CanSkip": true,
    "EndBranch": "end-branch",
    "EndTime": "0001-01-01T00:00:00Z"
  }
}`[1:]

		repoName := "git-town-unit-tests"
		err := runstate.Save(&runState, repoName)
		assert.NoError(t, err)
		filepath, err := runstate.PersistenceFilePath(repoName)
		assert.NoError(t, err)
		content, err := os.ReadFile(filepath)
		assert.NoError(t, err)
		assert.Equal(t, wantJSON, string(content))
		var newState runstate.RunState
		err = json.Unmarshal(content, &newState)
		assert.NoError(t, err)
		assert.Equal(t, runState, newState)
	})
}
