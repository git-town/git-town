package undo

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/domain"
)

// SingleCacheDiff describes the changes made to the local or global Git configuration.
type SingleCacheDiff struct {
	Added   []configdomain.Key
	Removed map[configdomain.Key]string
	Changed map[configdomain.Key]domain.Change[string]
}

func EmptySingleCacheDiff() SingleCacheDiff {
	return SingleCacheDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]domain.Change[string]{},
	}
}

func NewSingleCacheDiff(before, after gitconfig.SingleCache) SingleCacheDiff {
	result := SingleCacheDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]domain.Change[string]{},
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
