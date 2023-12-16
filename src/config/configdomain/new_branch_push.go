package configdomain

import "strconv"

// NewBranchPush indicates whether newly created branches should be pushed to the remote or not.
type NewBranchPush bool

func (self NewBranchPush) Bool() bool {
	return bool(self)
}

func (self NewBranchPush) String() string {
	return strconv.FormatBool(self.Bool())
}
