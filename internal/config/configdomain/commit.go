package configdomain

// when creating a new branch, whether to commit the currently staged changes into that new branch
type Commit bool

func (self Commit) IsTrue() bool {
	return bool(self)
}

func (self Commit) IsFalse() bool {
	return !self.IsTrue()
}
