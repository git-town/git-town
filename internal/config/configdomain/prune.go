package configdomain

// indicates whether a Git Town command should create a prototype branch
type Prune bool

func (self Prune) IsTrue() bool {
	return bool(self)
}
