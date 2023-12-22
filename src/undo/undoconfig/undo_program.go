package undoconfig

import (
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func DetermineUndoConfigProgram(initialConfigSnapshot undodomain.ConfigSnapshot, configGit *gitconfig.Access) (program.Program, error) {
	fullCache, err := gitconfig.LoadFullCache(configGit)
	if err != nil {
		return program.Program{}, err
	}
	finalConfigSnapshot := undodomain.ConfigSnapshot{
		GitConfig: fullCache,
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram(), nil
}
