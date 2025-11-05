package configdomain

// Commit indicates whether to commit the currently staged changes into a new branch.
type Commit bool

func (self Commit) ShouldCommit() bool {
	return bool(self)
}
