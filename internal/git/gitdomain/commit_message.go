package gitdomain

import "strings"

// CommitMessage is the entire textual messages of a Git commit.
type CommitMessage string

// CommitMessageParts describes the parts of a Git commit message.
type CommitMessageParts struct {
	Title CommitTitle // the first line of the commit message
	Body  string      // the commit message text minus the first line and empty lines separating it from the rest of the message
}

// Parts separates the parts of the given commit message.
func (self CommitMessage) Parts() CommitMessageParts {
	title, body, _ := strings.Cut(self.String(), "\n")
	for strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return CommitMessageParts{
		Title: CommitTitle(title),
		Body:  body,
	}
}

// String implements the fmt.Stringer interface.
func (self CommitMessage) String() string {
	return string(self)
}
