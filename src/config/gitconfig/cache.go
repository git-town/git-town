package gitconfig

import (
	"regexp"
	"sort"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"golang.org/x/exp/maps"
)

// Cache caches a Git configuration type (local or global).
type Cache map[configdomain.Key]string

// Clone provides a copy of this GitConfiguration instance.
func (self Cache) Clone() Cache {
	result := Cache{}
	maps.Copy(result, self)
	return result
}

// KeysMatching provides the keys in this GitConfigCache that match the given regex.
func (self Cache) KeysMatching(pattern string) []configdomain.Key {
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
