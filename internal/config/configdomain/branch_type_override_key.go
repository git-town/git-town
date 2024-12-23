package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

const BranchTypeSuffix = ".branchtype"

// a Key that contains a BranchTypeOverrides entry
type BranchTypeOverrideKey Key

func NewBranchTypeOverrideKey(key Key) Option[LineageKey] {
	if isBranchTypeOverrideKey(key.String()) {
		return Some(LineageKey(key))
	}
	return None[LineageKey]()
}

func (self LineageKey) BranchType() string {
	return strings.TrimSuffix(strings.TrimPrefix(self.String(), BranchSpecificKeyPrefix), BranchTypeSuffix)
}

func isBranchTypeOverrideKey(key string) bool {
	return strings.HasPrefix(key, BranchSpecificKeyPrefix) && strings.HasSuffix(key, BranchTypeSuffix)
}
