package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/vm/program"
)

func DetermineUndoConfigProgram(begin BeginConfigSnapshot, end EndConfigSnapshot) program.Program {
	configDiff := NewConfigDiffs(begin, end)
	return configDiff.UndoProgram()
}
