package undoconfig

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
)

// SingleCacheDiff provides a diff of the two given SingleCache instances.
func SingleCacheDiff(before, after configdomain.SingleSnapshot) ConfigDiff {
	result := ConfigDiff{
		Added:   []configdomain.Key{},
		Changed: map[configdomain.Key]undodomain.Change[string]{},
		Removed: map[configdomain.Key]string{},
	}
	for key, beforeValue := range before { // okay to iterate the map in random order because we assign to new maps
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
	for key := range mapstools.SortedKeys(after) {
		_, beforeContains := before[key]
		if !beforeContains {
			result.Added = append(result.Added, key)
		}
	}
	return result
}
