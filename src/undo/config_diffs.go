package undo

import (
	"github.com/git-town/git-town/v9/src/runstate"
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

func (cds ConfigDiffs) UndoSteps() runstate.StepList {
	result := runstate.StepList{}
	for _, key := range cds.Global.Added {
		result.Append(&steps.RemoveGlobalConfigStep{Key: key})
	}
	for key, value := range cds.Global.Removed {
		result.Append(&steps.SetGlobalConfigStep{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cds.Global.Changed {
		result.Append(&steps.SetGlobalConfigStep{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range cds.Local.Added {
		result.Append(&steps.RemoveLocalConfigStep{Key: key})
	}
	for key, value := range cds.Local.Removed {
		result.Append(&steps.SetLocalConfigStep{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cds.Local.Changed {
		result.Append(&steps.SetLocalConfigStep{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
