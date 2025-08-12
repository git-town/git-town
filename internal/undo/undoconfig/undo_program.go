package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/vm/program"
)

func DetermineUndoConfigProgram(begin BeginConfigSnapshot, end EndConfigSnapshot) program.Program {
	return NewConfigDiffs(begin, end).UndoProgram()
}
