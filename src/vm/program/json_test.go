package program_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		jsonstep := program.JSON{
			Opcode: &opcodes.Checkout{
				Branch: gitdomain.NewLocalBranchName("branch-1"),
			},
		}
		have, err := json.MarshalIndent(jsonstep, "", "  ")
		must.NoError(t, err)
		// NOTE: It's unclear why this doesn't contain the "data" and "type" fields from JSONStep's MarshalJSON method here.
		//       Marshaling an entire RunState somehow works correctly.
		want := `
{
  "Opcode": {
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
	"type": "Checkout"
}`[1:]
		have := program.JSON{
			Opcode: &opcodes.Checkout{
				Branch: "",
			},
		}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := program.JSON{
			Opcode: &opcodes.Checkout{
				Branch: gitdomain.NewLocalBranchName("branch-1"),
			},
		}
		must.Eq(t, want, have)
	})
}
