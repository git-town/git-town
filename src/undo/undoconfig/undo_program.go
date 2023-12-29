package undoconfig

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func DetermineUndoConfigProgram(initialConfigSnapshot ConfigSnapshot, configGit *configdomain.Access) (program.Program, error) {
	globalCache, _, err := configGit.LoadCache(true)
	if err != nil {
		return program.Program{}, err
	}
	localCache, _, err := configGit.LoadCache(false)
	if err != nil {
		return program.Program{}, err
	}
	finalConfigSnapshot := ConfigSnapshot{
		Global: globalCache,
		Local:  localCache,
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram(), nil
}
