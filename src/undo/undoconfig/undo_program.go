package undoconfig

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func DetermineUndoConfigProgram(initialConfigSnapshot configdomain.ConfigSnapshot, configGit *configdomain.Access) (program.Program, error) {
	fullCache, err := configdomain.LoadFullCache(configGit)
	if err != nil {
		return program.Program{}, err
	}
	finalConfigSnapshot := configdomain.ConfigSnapshot{
		GitConfig: fullCache,
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram(), nil
}
