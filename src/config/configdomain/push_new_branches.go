package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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

func ParsePushNewBranches(value, source string) (PushNewBranches, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	return PushNewBranches(parsed), nil
}

func ParsePushNewBranchesOption(value, source string) (Option[PushNewBranches], error) {
	result, err := ParsePushNewBranches(value, source)
	if err != nil {
		return None[PushNewBranches](), err
	}
	return Some(result), err
}
