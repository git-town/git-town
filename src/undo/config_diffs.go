package undo

import (
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// ConfigDiffs describes the changes made to the local and global Git configuration.
type ConfigDiffs struct {
	Global SingleCacheDiff
	Local  SingleCacheDiff
}

func NewConfigDiffs(before, after ConfigSnapshot) ConfigDiffs {
	return ConfigDiffs{
		Global: NewSingleCacheDiff(before.GitConfig.GlobalCache, after.GitConfig.GlobalCache),
		Local:  NewSingleCacheDiff(before.GitConfig.LocalCache, after.GitConfig.LocalCache),
	}
}

func (self ConfigDiffs) UndoProgram() program.Program {
	result := program.Program{}
	for _, key := range self.Global.Added {
		result.Add(&opcode.RemoveGlobalConfig{Key: key})
	}
	for key, value := range self.Global.Removed {
		result.Add(&opcode.SetGlobalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range self.Global.Changed {
		result.Add(&opcode.SetGlobalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range self.Local.Added {
		result.Add(&opcode.RemoveLocalConfig{Key: key})
	}
	for key, value := range self.Local.Removed {
		result.Add(&opcode.SetLocalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range self.Local.Changed {
		result.Add(&opcode.SetLocalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
