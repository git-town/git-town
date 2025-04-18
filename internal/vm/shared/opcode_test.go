package shared_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/vm/opcodes"
	"github.com/git-town/git-town/v19/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestOpcode(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.BranchTypeOverrideSet{Branch: "branch", BranchType: configdomain.BranchTypePerennialBranch},
		}
		have := give.String()
		want := `
Program:
1: &opcodes.MergeAbort{undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
2: &opcodes.BranchTypeOverrideSet{Branch:"branch", BranchType:"perennial", undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `
[
	{
		"data": {
			"Hard": false,
			"MustHaveSHA": "abcdef",
			"SetToSHA": "123456"
		},
		"type": "BranchCurrentResetToSHAIfNeeded"
	},
	{
		"data": {},
		"type": "StashOpenChanges"
	}
]`[1:]
		have := program.Program{}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := program.Program{
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				Hard:        false,
				MustHaveSHA: "abcdef",
				SetToSHA:    "123456",
			},
			&opcodes.StashOpenChanges{},
		}
		must.Eq(t, want, have)
	})
}
