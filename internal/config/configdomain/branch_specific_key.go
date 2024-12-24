package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
)

// a Git config key that contains a branch specific value
type BranchSpecificKey struct {
	Key
}

// provides the name of the child branch encoded in this LineageKey
func (self BranchSpecificKey) BranchName() gitdomain.LocalBranchName {
	text := strings.TrimSuffix(strings.TrimPrefix(self.String(), BranchSpecificKeyPrefix), LineageKeySuffix)
	return gitdomain.NewLocalBranchName(text)
}

const BranchSpecificKeyPrefix = "git-town-branch."

func isBranchSpecificKey(key string) bool {
	return strings.HasPrefix(key, BranchSpecificKeyPrefix)
}
