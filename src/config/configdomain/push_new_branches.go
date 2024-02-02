package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
)

// PushNewBranches indicates whether newly created branches should be pushed to the remote or not.
type PushNewBranches bool

func (self PushNewBranches) Bool() bool {
	return bool(self)
}

func (self PushNewBranches) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewNewBranchPushRef(value bool) *PushNewBranches {
	result := PushNewBranches(value)
	return &result
}

func ParseNewBranchPush(value, source string) (PushNewBranches, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	return PushNewBranches(parsed), nil
}

func ParseNewBranchPushRef(value, source string) (*PushNewBranches, error) {
	result, err := ParseNewBranchPush(value, source)
	return &result, err
}
