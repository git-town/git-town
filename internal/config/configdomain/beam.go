package configdomain

// indicates whether a Git Town command should create a prototype branch
type Beam bool

func (self Beam) ShouldBeam() bool {
	return bool(self)
}
