package configdomain

import "strconv"

// PushBranches indicates whether Git Town commands should push local commits to the respective tracking branch
type PushBranches bool

func (self PushBranches) ShouldPush() bool {
	return bool(self)
}

func (self PushBranches) String() string {
	return strconv.FormatBool(bool(self))
}
