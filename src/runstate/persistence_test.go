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
	t.Run("Serialization format", func(t *testing.T) {
		t.Parallel()
		t.Run("PushBranchStep", func(t *testing.T) {
			t.Parallel()
			step := steps.PushBranchStep{
				Branch:         domain.NewLocalBranchName("main"),
				ForceWithLease: true,
				NoPushHook:     true,
				Undoable:       true,
			}
			want := `
{
  "Branch": "main",
  "ForceWithLease": true,
  "NoPushHook": true,
  "Undoable": true
}`[1:]
			bytes, err := json.MarshalIndent(step, "", "  ")
			assert.NoError(t, err)
			assert.Equal(t, want, string(bytes))
			newStep := steps.PushBranchStep{}
			err = json.Unmarshal(bytes, &newStep)
			assert.NoError(t, err)
			assert.Equal(t, step, newStep)
		})
	})
}
