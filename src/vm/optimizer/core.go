// Package optimizer provides functionality to improve the performance of VM programs
// by re-arranging their opcodes.
// This is similar to optimizers in compilers.
package optimizer

import (
	"github.com/git-town/git-town/v14/src/vm/program"
)

// Optimize improves the performance of the given program by re-arranging its opcodes.
// It doesn't change the behavior of the program.
// This is similar to optimizers in compilers.
func Optimize(prog program.Program) program.Program {
	return RemoveDuplicateCheckout(prog)
}
