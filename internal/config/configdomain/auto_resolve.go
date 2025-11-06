package configdomain

// AutoResolve indicates whether a Git Town command should not auto-resolve phantom merge conflicts.
type AutoResolve bool

func (self AutoResolve) NoAutoResolve() bool {
	return !self.ShouldAutoResolve()
}

func (self AutoResolve) ShouldAutoResolve() bool {
	return bool(self)
}
