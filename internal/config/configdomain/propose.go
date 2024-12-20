package configdomain

// indicates whether a Git Town command should propose the branch it creates
type Propose bool

func (self Propose) IsTrue() bool {
	return bool(self)
}
