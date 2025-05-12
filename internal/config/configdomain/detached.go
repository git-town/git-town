package configdomain

// indicates whether a Git Town command should not update the root branch of the stack
type Detached bool

func (self Detached) IsFalse() bool {
	return !self.IsTrue()
}

func (self Detached) IsTrue() bool {
	return bool(self)
}
