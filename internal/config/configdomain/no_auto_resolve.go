package configdomain

// indicates whether a Git Town command should not auto-resolve phantom merge conflicts
type NoAutoResolve bool

func (self NoAutoResolve) IsTrue() bool {
	return bool(self)
}

func (self NoAutoResolve) ShouldAutoResolve() bool {
	return !self.IsTrue()
}
