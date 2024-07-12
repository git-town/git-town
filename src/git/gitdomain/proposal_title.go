package gitdomain

// CommitMessage is the entire textual messages of a Git commit.
type ProposalTitle string

// String implements the fmt.Stringer interface.
func (self ProposalTitle) String() string {
	return string(self)
}
