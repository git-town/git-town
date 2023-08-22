package runstate_test

import (
	"encoding/json"
	"reflect"
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
	t.Run("Serialization format", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			step steps.Step
			json string
		}{
			{
				step: &steps.PushBranchStep{
					Branch:         domain.NewLocalBranchName("main"),
					ForceWithLease: true,
					NoPushHook:     true,
					Undoable:       true,
				},
				json: `
{
  "Branch": "main",
  "ForceWithLease": true,
  "NoPushHook": true,
  "Undoable": true
}`[1:],
			},
		}
		for _, test := range tests {
			bytes, err := json.MarshalIndent(test.step, "", "  ")
			assert.NoError(t, err)
			assert.Equal(t, test.json, string(bytes))
			stepTypeName := "*" + reflect.TypeOf(test.step).String()[7:]
			newStep := runstate.DetermineStep(stepTypeName)
			err = json.Unmarshal(bytes, &newStep)
			assert.NoError(t, err)
			assert.Equal(t, test.step, newStep)
		}
	})
}
