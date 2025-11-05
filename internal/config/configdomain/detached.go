package configdomain

import "strconv"

// Detached indicates whether a Git Town command should not update the root branch of the stack.
type Detached bool

func (self Detached) ShouldWorkDetached() bool {
	return bool(self)
}

func (self Detached) String() string {
	return strconv.FormatBool(bool(self))
}
