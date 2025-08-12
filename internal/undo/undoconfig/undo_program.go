package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/vm/program"
)

func DetermineUndoConfigProgram(initialConfigSnapshot BeginConfigSnapshot, finalConfigSnapshot EndConfigSnapshot) program.Program {
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram()
}
