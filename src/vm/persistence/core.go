// Package persistence stores Git Town runstate on disk.
package persistence

import "github.com/git-town/git-town/v9/src/vm/shared"

type Opcode interface {
	CreateAbortProgram() []Opcode
	CreateContinueProgram() []Opcode
	CreateAutomaticAbortError() error
	Run(args shared.RunArgs) error
	ShouldAutomaticallyAbortOnError() bool
}
