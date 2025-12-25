package configdomain

import "strconv"

// AutoSync indicates whether a Git Town command should sync branches
// before performing its actual functionality.
type AutoSync bool

func (self AutoSync) ShouldSync() bool {
	return bool(self)
}

func (self AutoSync) String() string {
	return strconv.FormatBool(self.ShouldSync())
}
