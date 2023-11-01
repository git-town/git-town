package undo

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
)

// ConfigDiff describes the changes made to the local or global Git configuration.
type ConfigDiff struct {
	Added   []config.Key
	Removed map[config.Key]string
	Changed map[config.Key]domain.Change[string]
}

func EmptyConfigDiff() ConfigDiff {
	return ConfigDiff{
		Added:   []config.Key{},
		Removed: map[config.Key]string{},
		Changed: map[config.Key]domain.Change[string]{},
	}
}

func NewConfigDiff(before, after config.GitConfigCache) ConfigDiff {
	result := ConfigDiff{
		Added:   []config.Key{},
		Removed: map[config.Key]string{},
		Changed: map[config.Key]domain.Change[string]{},
	}
	for key, beforeValue := range before {
		afterValue, afterContains := after[key]
		if afterContains {
			if beforeValue != afterValue {
				result.Changed[key] = domain.Change[string]{
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
