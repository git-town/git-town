package shared

import (
	"github.com/git-town/git-town/v12/src/gohacks"
)

func IsEndOfBranchProgramOpcode(opcode Opcode) bool {
	typeName := gohacks.TypeName(opcode)
	return typeName == "EndOfBranchProgram"
}
