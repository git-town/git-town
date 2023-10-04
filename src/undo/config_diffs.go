package undo

import (
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
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

func (cds ConfigDiffs) UndoSteps() steps.List {
	result := steps.List{}
	for _, key := range cds.Global.Added {
		result.Append(&step.RemoveGlobalConfig{Key: key})
	}
	for key, value := range cds.Global.Removed {
		result.Append(&step.SetGlobalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cds.Global.Changed {
		result.Append(&step.SetGlobalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range cds.Local.Added {
		result.Append(&step.RemoveLocalConfig{Key: key})
	}
	for key, value := range cds.Local.Removed {
		result.Append(&step.SetLocalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cds.Local.Changed {
		result.Append(&step.SetLocalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
