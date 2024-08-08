package configdomain

// indicates whether to sync all branches or only the current branch
type FullStack bool

func (self FullStack) Enabled() bool {
	return bool(self)
}
