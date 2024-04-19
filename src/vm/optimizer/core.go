// Package optimizer provides functionality to improve the performance of VM programs
// by re-arranging their opcodes.
// This is similar to optimizers in compilers.
package optimizer

import (
	"github.com/git-town/git-town/v14/src/vm/program"
)

func Optimize(prog program.Program) program.Program {
	return RemoveDuplicateCheckout(prog)
}
