package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsEndOfBranchProgramOpcode(t *testing.T) {
	t.Parallel()
	tests := map[shared.Opcode]bool{
		&opcodes.ProgramEndOfBranch{}: true,
		&opcodes.MergeAbort{}:         false,
	}
	for give, want := range tests {
		have := shared.IsEndOfBranchProgramOpcode(give)
		must.Eq(t, want, have)
	}
}
