package gitconfig

import (
	"regexp"
	"sort"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"golang.org/x/exp/maps"
)

// SingleCache caches a single Git configuration type (local or global).
type SingleCache map[configdomain.Key]string

// Clone provides a copy of this GitConfiguration instance.
func (self SingleCache) Clone() SingleCache {
	result := SingleCache{}
	maps.Copy(result, self)
	return result
}

// KeysMatching provides the keys in this GitConfigCache that match the given regex.
func (self SingleCache) KeysMatching(pattern string) []configdomain.Key {
	result := []configdomain.Key{}
	re := regexp.MustCompile(pattern)
	for key := range self {
		if re.MatchString(key.String()) {
			result = append(result, key)
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].String() < result[b].String() })
	return result
}

// SingleCacheDiff provides a diff of the two given SingleCache instances.
func SingleCacheDiff(before, after SingleCache) configdomain.ConfigDiff {
	result := configdomain.ConfigDiff{
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
