package shared_test

import (
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
}
