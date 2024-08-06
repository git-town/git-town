package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsEndOfBranchProgramOpcode(t *testing.T) {
	t.Parallel()
	tests := map[shared.Opcode]bool{
		&opcodes.EndOfBranchProgram{}: true,
		&opcodes.AbortMerge{}:         false,
	}
	for give, want := range tests {
		have := shared.IsEndOfBranchProgramOpcode(give)
		must.Eq(t, want, have)
	}
}
