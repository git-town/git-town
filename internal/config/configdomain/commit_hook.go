package configdomain

import (
	"strconv"
)

// CommitHook indicates whether commit-hooks are enabled.
type CommitHook bool

func (self CommitHook) String() string {
	return strconv.FormatBool(bool(self))
}

const (
	CommitHookEnabled  CommitHook = true
	CommitHookDisabled CommitHook = false
)
