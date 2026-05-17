package gitdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// CommitMessageFile is the name of the file from which to read the CommitMessage.
type CommitMessageFile stringss.TrimmedString

// ShouldReadStdin indicates whether the commit message should be read from STDIN.
func (self CommitMessageFile) ShouldReadStdin() bool {
	return self == "-"
}

// String implements the fmt.Stringer interface.
func (self CommitMessageFile) String() string {
	return string(self)
}
