package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/internal/gohacks/prelude"
)

// a Key that contains a lineage entry
type LineageKey Key

// CheckLineage indicates using the returned option whether this key is a lineage key.
func NewLineageKey(key Key) Option[LineageKey] {
	if isLineageKey(key.String()) {
		return Some(LineageKey(key))
	}
	return None[LineageKey]()
}

// provides the name of the child branch encoded in this LineageKey
func (self LineageKey) ChildName() string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(self.String(), LineageKeyPrefix), LineageKeySuffix))
}

// converts this LineageKey into a generic Key
func (self LineageKey) Key() Key {
	return Key(self)
}

func (self LineageKey) String() string {
	return string(self)
}

const (
	LineageKeyPrefix = "git-town-branch."
	LineageKeySuffix = ".parent"
)

// indicates whether the given key value is for a LineageKey
func isLineageKey(key string) bool {
	return strings.HasPrefix(key, LineageKeyPrefix) && strings.HasSuffix(key, LineageKeySuffix)
}
