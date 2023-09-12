package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/steps"
)

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Cwd       string // the current working directory
	GitConfig config.GitConfig
}

func (cs ConfigSnapshot) Diff(other ConfigSnapshot) SnapshotConfigDiff {
	return SnapshotConfigDiff{
		Global: diffConfig(cs.GitConfig.Global, other.GitConfig.Global),
		Local:  diffConfig(cs.GitConfig.Local, other.GitConfig.Local),
	}
}

func diffConfig(before, after config.GitConfigCache) ConfigDiff {
	result := ConfigDiff{
		Added:   []config.Key{},
		Removed: map[config.Key]string{},
		Changed: map[config.Key]Change[string]{},
	}
	for key, beforeValue := range before {
		afterValue, afterContains := after[key]
		if afterContains {
			if beforeValue != afterValue {
				result.Changed[key] = Change[string]{
					Before: beforeValue,
					After:  afterValue,
				}
			}
		} else {
			result.Removed[key] = beforeValue
		}
	}
	for key := range after {
		_, beforeContains := before[key]
		if !beforeContains {
			result.Added = append(result.Added, key)
		}
	}
	return result
}

type ConfigDiff struct {
	Added   []config.Key
	Removed map[config.Key]string
	Changed map[config.Key]Change[string]
}

type SnapshotConfigDiff struct {
	Global ConfigDiff
	Local  ConfigDiff
}

func (scd SnapshotConfigDiff) UndoSteps() StepList {
	result := StepList{}
	for _, key := range scd.Global.Added {
		result.Append(&steps.RemoveGlobalConfigStep{Key: key})
	}
	for _, key := range scd.Local.Added {
		result.Append(&steps.RemoveLocalConfigStep{Key: key})
	}
	for key, value := range scd.Global.Removed {
		result.Append(&steps.SetGlobalConfigStep{
			Key:   key,
			Value: value,
		})
	}
	for key, value := range scd.Local.Removed {
		result.Append(&steps.SetLocalConfigStep{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range scd.Global.Changed {
		result.Append(&steps.SetGlobalConfigStep{
			Key:   key,
			Value: change.Before,
		})
	}
	for key, change := range scd.Local.Changed {
		result.Append(&steps.SetLocalConfigStep{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}

type Change[T any] struct {
	Before T
	After  T
}
