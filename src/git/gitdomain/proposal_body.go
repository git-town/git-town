package gitdomain

// CommitMessage is the entire textual messages of a Git commit.
type ProposalBody string

// String implements the fmt.Stringer interface.
func (self ProposalBody) String() string {
	return string(self)
}
