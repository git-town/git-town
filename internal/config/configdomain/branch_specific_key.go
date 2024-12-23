package configdomain

import "strings"

// a Git config key that contains a branch specific value
type BranchSpecificKey struct {
	Key
}

// provides the name of the child branch encoded in this LineageKey
func (self BranchSpecificKey) BranchName() string {
	return strings.TrimSuffix(strings.TrimPrefix(self.String(), BranchSpecificKeyPrefix), LineageKeySuffix)
}

const BranchSpecificKeyPrefix = "git-town-branch."

func isBranchSpecificKey(key string) bool {
	return strings.HasPrefix(key, BranchSpecificKeyPrefix)
}
