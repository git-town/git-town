package configdomain

// indicates whether a Git Town command should create a prototype branch
type Prototype bool

func (self Prototype) IsFalse() bool {
	return !bool(self)
}

func (self Prototype) IsTrue() bool {
	return bool(self)
}
