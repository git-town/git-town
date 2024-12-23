package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// a Key that contains a lineage entry
type LineageKey struct {
	Key
}

func NewLineageKey(key Key) LineageKey {
	return LineageKey{
		Key: key,
	}
}

// CheckLineage indicates using the returned option whether this key is a lineage key.
func ParseLineageKey(key Key) Option[LineageKey] {
	if isLineageKey(key.String()) {
		return Some(LineageKey{
			Key: key,
		})
	}
	return None[LineageKey]()
}

// provides the name of the child branch encoded in this LineageKey
func (self LineageKey) ChildName() string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(self.String(), LineageKeyPrefix), LineageKeySuffix))
}

const (
	LineageKeyPrefix = "git-town-branch."
	LineageKeySuffix = ".parent"
)

// indicates whether the given key value is for a LineageKey
func isLineageKey(key string) bool {
	return strings.HasPrefix(key, LineageKeyPrefix) && strings.HasSuffix(key, LineageKeySuffix)
}
