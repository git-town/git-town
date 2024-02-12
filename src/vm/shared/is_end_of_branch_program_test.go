package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/vm/opcode"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsEndOfBranchProgramOpcode(t *testing.T) {
	t.Parallel()

	t.Run("given an opcode.EndOfBranchProgram", func(t *testing.T) {
		t.Parallel()
		give := &opcode.EndOfBranchProgram{}
		must.True(t, shared.IsEndOfBranchProgramOpcode(give))
	})

	t.Run("given another opcode", func(t *testing.T) {
		t.Parallel()
		give := &opcode.AbortMerge{}
		must.False(t, shared.IsEndOfBranchProgramOpcode(give))
	})
}
