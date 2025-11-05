package gitdomain

// CommitTitle is the first line of a CommitMessage.
type CommitTitle string

func (self CommitTitle) String() string {
	return string(self)
}
