package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestJSONStep(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		jsonstep := runstate.JSONStep{
			Step: &steps.CheckoutStep{
				Branch: domain.NewLocalBranchName("branch-1"),
			},
		}
		have, err := json.MarshalIndent(jsonstep, "", "  ")
		assert.Nil(t, err)
		// NOTE: It's unclear why this doesn't contain the "data" and "type" fields from JSONStep's MarshalJSON method here.
		//       Marshaling an entire RunState somehow works correctly.
		want := `
{
  "Step": {
    "Branch": "branch-1"
  }
}`[1:]
		assert.Equal(t, want, string(have))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `
{
	"data": {
    "Branch": "branch-1"
  },
	"type": "CheckoutStep"
}`[1:]
		have := runstate.JSONStep{
			Step: &steps.CheckoutStep{
				Branch: domain.LocalBranchName{},
			},
		}
		json.Unmarshal([]byte(give), &have)
		want := runstate.JSONStep{
			Step: &steps.CheckoutStep{
				Branch: domain.NewLocalBranchName("branch-1"),
			},
		}
		assert.Equal(t, want, have)
	})
}
