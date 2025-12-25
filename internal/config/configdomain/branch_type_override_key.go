package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const BranchTypeSuffix = ".branchtype"

// BranchTypeOverrideKey is a Key that contains a BranchTypeOverrides entry,
// for example "git-town-branch.foo.branchtype".
type BranchTypeOverrideKey struct {
	BranchSpecificKey
}

// Branch provides the name of the child branch encoded in this LineageKey.
func (self BranchTypeOverrideKey) Branch() gitdomain.LocalBranchName {
	text := strings.TrimSuffix(strings.TrimPrefix(self.String(), BranchSpecificKeyPrefix), BranchTypeSuffix)
	return gitdomain.NewLocalBranchName(text)
}

func IsBranchTypeOverrideKey(key string) bool {
	return isBranchSpecificKey(key) && strings.HasSuffix(key, BranchTypeSuffix)
}

func NewBranchTypeOverrideKeyForBranch(branch gitdomain.LocalBranchName) BranchTypeOverrideKey {
	return BranchTypeOverrideKey{
		BranchSpecificKey: BranchSpecificKey{
			Key: Key(BranchSpecificKeyPrefix + branch.String() + BranchTypeSuffix),
		},
	}
}

func ParseBranchTypeOverrideKey(key Key) Option[BranchTypeOverrideKey] {
	if IsBranchTypeOverrideKey(key.String()) {
		return Some(BranchTypeOverrideKey{
			BranchSpecificKey: BranchSpecificKey{
				Key: key,
			},
		})
	}
	return None[BranchTypeOverrideKey]()
}
