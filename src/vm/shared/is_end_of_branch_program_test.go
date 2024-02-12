package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/vm/opcode"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsEndOfBranchProgramOpcode(t *testing.T) {
	t.Parallel()
	tests := map[shared.Opcode]bool{
		&opcode.EndOfBranchProgram{}: true,
		&opcode.AbortMerge{}:         false,
	}
	for give, want := range tests {
		have := shared.IsEndOfBranchProgramOpcode(give)
		must.Eq(t, want, have)
	}
}
