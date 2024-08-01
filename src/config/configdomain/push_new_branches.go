package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// PushNewBranches indicates whether newly created branches should be pushed to the remote or not.
type PushNewBranches bool

func (self PushNewBranches) Bool() bool {
	return bool(self)
}

func (self PushNewBranches) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParsePushNewBranches(value, source string) (Option[PushNewBranches], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(PushNewBranches(parsed)), err
	}
	return None[PushNewBranches](), err
}
