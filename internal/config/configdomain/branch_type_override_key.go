package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

const BranchTypeSuffix = ".branchtype"

// a Key that contains a BranchTypeOverrides entry,
// for example "git-town-branch.foo.branchtype"
type BranchTypeOverrideKey struct {
	BranchSpecificKey
}

func ParseBranchTypeOverrideKey(key Key) Option[BranchTypeOverrideKey] {
	if isBranchTypeOverrideKey(key.String()) {
		return Some(BranchTypeOverrideKey{
			BranchSpecificKey: BranchSpecificKey{
				Key: key,
			},
		})
	}
	return None[BranchTypeOverrideKey]()
}

func (self BranchTypeOverrideKey) Key() Key {
	return Key(self.BranchSpecificKey.Key)
}

func isBranchTypeOverrideKey(key string) bool {
	return strings.HasPrefix(key, BranchSpecificKeyPrefix) && strings.HasSuffix(key, BranchTypeSuffix)
}
