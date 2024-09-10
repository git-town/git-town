package configdomain

// indicates whether to sync all branches or only the current branch
type AllBranches bool

func (self AllBranches) Enabled() bool {
	return bool(self)
}
