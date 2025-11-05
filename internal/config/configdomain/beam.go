package configdomain

// Beam indicates whether a Git Town command should beam commits to the new branch.
type Beam bool

func (self Beam) ShouldBeam() bool {
	return bool(self)
}
