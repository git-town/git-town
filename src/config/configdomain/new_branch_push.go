package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
)

// NewBranchPush indicates whether newly created branches should be pushed to the remote or not.
type NewBranchPush bool

func (self NewBranchPush) Bool() bool {
	return bool(self)
}

func (self NewBranchPush) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewNewBranchPushRef(value bool) *NewBranchPush {
	result := NewBranchPush(value)
	return &result
}

func ParseNewBranchPushRef(value string) (*NewBranchPush, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(messages.ValueInvalid, KeyPushNewBranches, value)
	}
	token := NewBranchPush(parsed)
	return &token, nil
}
