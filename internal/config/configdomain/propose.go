package configdomain

// indicates whether a Git Town command should create a prototype branch
type Propose bool

func (self Propose) IsTrue() bool {
	return bool(self)
}
