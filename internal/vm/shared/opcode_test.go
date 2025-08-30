package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	"github.com/shoenig/test/must"
)

func TestOpcode(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := []shared.Opcode{
			&opcodes.MergeAbort{},
			&opcodes.BranchTypeOverrideSet{Branch: "branch", BranchType: configdomain.BranchTypePerennialBranch},
		}
		have := shared.RenderOpcodes(give, "")
		want := `
Program:
1: &opcodes.MergeAbort{undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
2: &opcodes.BranchTypeOverrideSet{Branch:"branch", BranchType:"perennial", undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
`[1:]
		must.EqOp(t, want, have)
	})
}
