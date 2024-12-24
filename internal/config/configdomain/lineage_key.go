package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// a Key that contains a lineage entry
type LineageKey struct {
	BranchSpecificKey
}

func NewLineageKey(key Key) LineageKey {
	return LineageKey{
		BranchSpecificKey: BranchSpecificKey{
			Key: key,
		},
	}
}

// CheckLineage indicates using the returned option whether this key is a lineage key.
func ParseLineageKey(key Key) Option[LineageKey] {
	if isLineageKey(key.String()) {
		return Some(NewLineageKey(key))
	}
	return None[LineageKey]()
}

const LineageKeySuffix = ".parent"

// indicates whether the given key value is for a LineageKey
func isLineageKey(key string) bool {
	return isBranchSpecificKey(key) && strings.HasSuffix(key, LineageKeySuffix)
}
