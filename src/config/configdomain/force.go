package configdomain

// indicates whether a Git Town command should execute the commands or only display them
type Force bool

func (self Force) IsFalse() bool {
	return !bool(self)
}
