package gitdomain

// CommitMessageFile is the name of the file from which to read the CommitMessage.
type CommitMessageFile string

// ShouldReadStdin indicates whether the commit message should be read from STDIN.
func (self CommitMessageFile) ShouldReadStdin() bool {
	return self == "-"
}

// String implements the fmt.Stringer interface.
func (self CommitMessageFile) String() string {
	return string(self)
}
