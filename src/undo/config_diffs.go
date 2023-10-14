package undo

import (
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/git-town/git-town/v9/src/vm/step"
)

// ConfigDiffs describes the changes made to the local and global Git configuration.
type ConfigDiffs struct {
	Global ConfigDiff
	Local  ConfigDiff
}

func NewConfigDiffs(before, after ConfigSnapshot) ConfigDiffs {
	return ConfigDiffs{
		Global: NewConfigDiff(before.GitConfig.Global, after.GitConfig.Global),
		Local:  NewConfigDiff(before.GitConfig.Local, after.GitConfig.Local),
	}
}

func (cds ConfigDiffs) UndoProgram() program.Program {
	result := program.Program{}
	for _, key := range cds.Global.Added {
		result.Add(&step.RemoveGlobalConfig{Key: key})
	}
	for key, value := range cds.Global.Removed {
		result.Add(&step.SetGlobalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cds.Global.Changed {
		result.Add(&step.SetGlobalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range cds.Local.Added {
		result.Add(&step.RemoveLocalConfig{Key: key})
	}
	for key, value := range cds.Local.Removed {
		result.Add(&step.SetLocalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cds.Local.Changed {
		result.Add(&step.SetLocalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
