package configdomain

import "strconv"

// Stash indicates whether Git Town should stash open changes before creating a new branch
type Stash bool

func (self Stash) ShouldStash() bool {
	return bool(self)
}

func (self Stash) String() string {
	return strconv.FormatBool(self.ShouldStash())
}
