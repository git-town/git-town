package undoconfig

import (
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/undo/undodomain"
	"github.com/git-town/git-town/v14/pkg/keys"
)

// SingleCacheDiff provides a diff of the two given SingleCache instances.
func SingleCacheDiff(before, after configdomain.SingleSnapshot) ConfigDiff {
	result := ConfigDiff{
		Added:   []keys.Key{},
		Changed: map[keys.Key]undodomain.Change[string]{},
		Removed: map[keys.Key]string{},
	}
	for key, beforeValue := range before {
		afterValue, afterContains := after[key]
		if afterContains {
			if beforeValue != afterValue {
				result.Changed[key] = undodomain.Change[string]{
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
