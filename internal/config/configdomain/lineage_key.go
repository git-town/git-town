package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// a Key that contains a lineage entry
type LineageKey BranchSpecificKey

func NewLineageKey(key Key) Option[LineageKey] {
	if isLineageKey(key.String()) {
		return Some(LineageKey(key))
	}
	return None[LineageKey]()
}

// converts this LineageKey into a generic Key
func (self LineageKey) Key() Key {
	return Key(self)
}

func (self LineageKey) String() string {
	return string(self)
}

const (
	BranchSpecificKeyPrefix = "git-town-branch."
	LineageKeySuffix        = ".parent"
)

// indicates whether the given key value is for a LineageKey
func isLineageKey(key string) bool {
	return strings.HasPrefix(key, BranchSpecificKeyPrefix) && strings.HasSuffix(key, LineageKeySuffix)
}
