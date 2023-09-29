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

func (cd ConfigDiffs) UndoSteps() runstate.StepList {
	result := runstate.StepList{}
	for _, key := range cd.Global.Added {
		result.Append(&steps.RemoveGlobalConfigStep{Key: key})
	}
	for key, value := range cd.Global.Removed {
		result.Append(&steps.SetGlobalConfigStep{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cd.Global.Changed {
		result.Append(&steps.SetGlobalConfigStep{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range cd.Local.Added {
		result.Append(&steps.RemoveLocalConfigStep{Key: key})
	}
	for key, value := range cd.Local.Removed {
		result.Append(&steps.SetLocalConfigStep{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range cd.Local.Changed {
		result.Append(&steps.SetLocalConfigStep{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
