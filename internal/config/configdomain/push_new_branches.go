package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v18/internal/gohacks"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// PushNewBranches indicates whether newly created branches should be pushed to the remote or not.
type PushNewBranches bool

func (self PushNewBranches) IsTrue() bool {
	return bool(self)
}

func (self PushNewBranches) String() string {
	return strconv.FormatBool(bool(self))
}

func ParsePushNewBranches(value string, source Key) (Option[PushNewBranches], error) {
	parsedOpt, err := gohacks.ParseBool(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(PushNewBranches(parsed)), err
	}
	return None[PushNewBranches](), err
}
