package configdomain

import (
	"strings"
)

// a Git config key that contains a branch specific value
type BranchSpecificKey struct {
	Key
}

const BranchSpecificKeyPrefix = "git-town-branch."

func isBranchSpecificKey(key string) bool {
	return strings.HasPrefix(key, BranchSpecificKeyPrefix)
}
