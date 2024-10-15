package undoconfig

import (
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
)

// ConfigDiffs describes the changes made to the local and global Git configuration.
type ConfigDiffs struct {
	Global ConfigDiff
	Local  ConfigDiff
}

func NewConfigDiffs(before, after ConfigSnapshot) ConfigDiffs {
	return ConfigDiffs{
		Global: SingleCacheDiff(before.Global, after.Global),
		Local:  SingleCacheDiff(before.Local, after.Local),
	}
}

func (self ConfigDiffs) UndoProgram() program.Program {
	result := program.Program{}
	for _, key := range self.Global.Added {
		result.Add(&opcodes.ConfigGlobalRemove{Key: key})
	}
	for key, value := range self.Global.Removed {
		result.Add(&opcodes.ConfigGlobalSet{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range self.Global.Changed {
		result.Add(&opcodes.ConfigGlobalSet{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range self.Local.Added {
		result.Add(&opcodes.ConfigLocalRemove{Key: key})
	}
	for key, value := range self.Local.Removed {
		result.Add(&opcodes.SetLocalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range self.Local.Changed {
		result.Add(&opcodes.SetLocalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
