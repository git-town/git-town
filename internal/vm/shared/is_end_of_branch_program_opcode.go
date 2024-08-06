package shared

import (
	"github.com/git-town/git-town/v14/internal/gohacks"
)

func IsEndOfBranchProgramOpcode(opcode Opcode) bool {
	typeName := gohacks.TypeName(opcode)
	return typeName == "EndOfBranchProgram"
}
