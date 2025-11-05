package configdomain

// Beam indicates whether a Git Town command should beam commits onto the created branch.
type Beam bool

// ShouldBeam indicates whether the beaming is enabled.
func (self Beam) ShouldBeam() bool {
	return bool(self)
}
