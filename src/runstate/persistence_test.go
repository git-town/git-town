package runstate_test

import (
	"encoding/json"
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
					&steps.AddToPerennialBranchesStep{Branch: domain.NewLocalBranchName("branch")},
				},
			},
			UndoStepList:      runstate.StepList{},
			UnfinishedDetails: &runstate.UnfinishedRunStateDetails{},
		}
		bytes, err := json.MarshalIndent(runState, "", "  ")
		assert.NoError(t, err)
		wantJSON := `
{
  "List": [
    {},
    {
      "Branch": "branch"
    }
  ]
}`[1:]
		assert.Equal(t, wantJSON, string(bytes))
		newState := runstate.RunState{}
		err = json.Unmarshal(bytes, &newList)
		assert.NoError(t, err)
		assert.Equal(t, stepList, newList)
	})
}
