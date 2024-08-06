package configdomain

// indicates whether a Git Town command should execute the commands despite not all safety conditions in place
type Force bool

func (self Force) IsFalse() bool {
	return !bool(self)
}
