package configdomain

// indicates whether a Git Town command should create a prototype branch
type Beam bool

func (self Beam) IsFalse() bool {
	return !self.IsTrue()
}

func (self Beam) IsTrue() bool {
	return bool(self)
}
