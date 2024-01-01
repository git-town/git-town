package undoconfig

import (
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func DetermineUndoConfigProgram(initialConfigSnapshot ConfigSnapshot, configGit *gitconfig.Access) (program.Program, error) {
	globalCache, _, err := configGit.Load(true)
	if err != nil {
		return program.Program{}, err
	}
	localCache, _, err := configGit.Load(false)
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
