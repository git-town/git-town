package configdomain

import "strings"

// a Key that contains a lineage entry
type LineageKey Key

func (self LineageKey) ChildName() string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(self.String(), LineageKeyPrefix), LineageKeySuffix))
}

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

func isLineageKey(key string) bool {
	return strings.HasPrefix(key, LineageKeyPrefix) && strings.HasSuffix(key, LineageKeySuffix)
}
