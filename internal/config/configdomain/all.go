package configdomain

// indicates whether to sync all branches or only the current branch
type SyncAllBranches bool

func (self SyncAllBranches) Enabled() bool {
	return bool(self)
}
