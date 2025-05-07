package configdomain

import (
	"strconv"
)

// PushHook contains the push-hook configuration setting.
type CommitHook bool

const (
	CommitHookEnabled  CommitHook = true
	CommitHookDisabled CommitHook = false
)

func (self CommitHook) String() string {
	return strconv.FormatBool(bool(self))
}
