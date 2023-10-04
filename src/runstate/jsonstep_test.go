package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/shoenig/test/must"
)

func TestJSONStep(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		jsonstep := runstate.JSONStep{
			Step: &step.Checkout{
				Branch: domain.NewLocalBranchName("branch-1"),
			},
		}
		have, err := json.MarshalIndent(jsonstep, "", "  ")
		must.NoError(t, err)
		// NOTE: It's unclear why this doesn't contain the "data" and "type" fields from JSONStep's MarshalJSON method here.
		//       Marshaling an entire RunState somehow works correctly.
		want := `
{
  "Step": {
    "Branch": "branch-1"
  }
}`[1:]
		must.EqOp(t, want, string(have))
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
			Step: &step.Checkout{
				Branch: domain.EmptyLocalBranchName(),
			},
		}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := runstate.JSONStep{
			Step: &step.Checkout{
				Branch: domain.NewLocalBranchName("branch-1"),
			},
		}
		must.Eq(t, want, have)
	})
}
