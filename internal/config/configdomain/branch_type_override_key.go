package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

const BranchTypeSuffix = ".branchtype"

// a Key that contains a BranchTypeOverrides entry,
// for example "git-town-branch.foo.branchtype"
type BranchTypeOverrideKey struct {
	BranchSpecificKey
}

func NewBranchTypeOverrideKeyForBranch(branch gitdomain.LocalBranchName) BranchTypeOverrideKey {
	return BranchTypeOverrideKey{
		BranchSpecificKey: BranchSpecificKey{
			Key: Key(BranchSpecificKeyPrefix + branch.String() + BranchTypeSuffix),
		},
	}
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

// provides the name of the child branch encoded in this LineageKey
func (self BranchTypeOverrideKey) Branch() gitdomain.LocalBranchName {
	text := strings.TrimSuffix(strings.TrimPrefix(self.String(), BranchSpecificKeyPrefix), BranchTypeSuffix)
	return gitdomain.NewLocalBranchName(text)
}

func isBranchTypeOverrideKey(key string) bool {
	return isBranchSpecificKey(key) && strings.HasSuffix(key, BranchTypeSuffix)
}
