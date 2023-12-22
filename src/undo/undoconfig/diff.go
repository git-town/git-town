package undoconfig

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
)

// SingleCacheDiff provides a diff of the two given SingleCache instances.
func SingleCacheDiff(before, after gitconfig.SingleCache) undodomain.ConfigDiff {
	result := undodomain.ConfigDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]undodomain.Change[string]{},
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
