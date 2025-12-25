// Package opcodes defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcodes

import (
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

func IsCheckoutOpcode(opcode shared.Opcode) bool {
	switch opcode.(type) {
	case *Checkout, *CheckoutIfExists, *CheckoutIfNeeded:
		return true
	default:
		return false
	}
}

func IsEndOfBranchProgramOpcode(opcode shared.Opcode) bool {
	_, ok := opcode.(*ProgramEndOfBranch)
	return ok
}

func Lookup(opcodeType string) shared.Opcode { //nolint:ireturn
	for _, opcode := range All() {
		if gohacks.TypeName(opcode) == opcodeType {
			return opcode
		}
	}
	return nil
}
