package configdomain

// indicates whether a Git Town command should create a prototype branch
type Beam bool

func (self Beam) IsTrue() bool {
	return bool(self)
}
