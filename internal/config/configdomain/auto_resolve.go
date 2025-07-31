package configdomain

// indicates whether a Git Town command should not auto-resolve phantom merge conflicts
type AutoResolve bool

func (self AutoResolve) ShouldAutoResolve() bool {
	return bool(self)
}
