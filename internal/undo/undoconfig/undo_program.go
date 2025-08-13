package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/vm/program"
)

func DetermineUndoConfigProgram(begin configdomain.BeginConfigSnapshot, end configdomain.EndConfigSnapshot) program.Program {
	return NewConfigDiffs(begin, end).UndoProgram()
}
