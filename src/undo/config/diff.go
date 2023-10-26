package config

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
)

// Diff describes the changes made to the local or global Git configuration.
type Diff struct {
	Added   []config.Key
	Removed map[config.Key]string
	Changed map[config.Key]domain.Change[string]
}

func EmptyDiff() Diff {
	return Diff{
		Added:   []config.Key{},
		Removed: map[config.Key]string{},
		Changed: map[config.Key]domain.Change[string]{},
	}
}

func NewDiff(before, after config.GitConfigCache) Diff {
	result := Diff{
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
