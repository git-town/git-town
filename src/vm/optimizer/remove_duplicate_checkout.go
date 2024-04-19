package optimizer

import (
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RemoveDuplicateCheckout returns the given program were checkout opcodes
// that are immediately followed by other checkout opcodes are removed.
func RemoveDuplicateCheckout(prog program.Program) program.Program {
	result := make([]shared.Opcode, 0, len(prog))
	var lastOpcode shared.Opcode
	for _, opcode := range prog {
		if shared.IsCheckoutOpcode(opcode) {
			lastOpcode = opcode
			continue
		}
		if shared.IsEndOfBranchProgramOpcode(opcode) {
			result = append(result, opcode)
			continue
		}
		if lastOpcode != nil {
			result = append(result, lastOpcode)
		}
		lastOpcode = nil
		result = append(result, opcode)
	}
	if lastOpcode != nil {
		result = append(result, lastOpcode)
	}
	return result
}
