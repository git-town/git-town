package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/messages"
)

// PushNewBranches indicates whether newly created branches should be pushed to the remote or not.
type PushNewBranches bool

func (self PushNewBranches) Bool() bool {
	return bool(self)
}

func (self PushNewBranches) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewPushNewBranchesRef(value bool) *PushNewBranches {
	result := PushNewBranches(value)
	return &result
}

func ParsePushNewBranches(value, source string) (PushNewBranches, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	return PushNewBranches(parsed), nil
}

func ParsePushNewBranchesRef(value, source string) (*PushNewBranches, error) {
	result, err := ParsePushNewBranches(value, source)
	return &result, err
}
