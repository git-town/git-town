package runstate_test

import (
	"encoding/json"
	"os"
	"testing"

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
						StartingPoint: domain.Location(domain.NewSHA("123456").Location()),
					},
				},
			},
			UndoStepList:      runstate.StepList{},
			UnfinishedDetails: &runstate.UnfinishedRunStateDetails{},
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
    }
  ],
  "UndoStepList": [],
  "UnfinishedDetails": {
    "CanSkip": false,
    "EndBranch": "",
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
		haveJSON := string(content)
		assert.Equal(t, wantJSON, haveJSON)
		var newState runstate.RunState
		err = json.Unmarshal(content, &newState)
		assert.NoError(t, err)
		assert.Equal(t, runState, newState)
	})
}
