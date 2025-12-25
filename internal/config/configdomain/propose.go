package configdomain

// Propose indicates whether a Git Town command should propose the branch it creates
type Propose bool

func (self Propose) ShouldPropose() bool {
	return bool(self)
}
