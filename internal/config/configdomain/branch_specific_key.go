package configdomain

import "strings"

// a Git config key that contains a branch specific value
type BranchSpecificKey Key

// provides the name of the child branch encoded in this LineageKey
func (self BranchSpecificKey) ChildName() string {
	return strings.TrimSuffix(strings.TrimPrefix(self.String(), BranchSpecificKeyPrefix), LineageKeySuffix)
}

func (self BranchSpecificKey) String() string {
	return string(self)
}
