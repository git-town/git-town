package configdomain

import "strconv"

// PushHook contains the push-hook configuration setting.
type PushHook bool

func (self PushHook) ShouldRunPushHook() bool {
	return bool(self)
}

func (self PushHook) String() string {
	return strconv.FormatBool(bool(self))
}
