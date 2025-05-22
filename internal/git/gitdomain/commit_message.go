package gitdomain

import "strings"

// CommitMessage is the entire textual messages of a Git commit.
type CommitMessage string

// CommitMessageParts describes the parts of a Git commit message.
type CommitMessageParts struct {
	Subject string // the first line of the commit message
	Text    string // the commit message text minus the first line and empty lines separating it from the rest of the message
}

// Parts separates the parts of the given commit message.
func (self CommitMessage) Parts() CommitMessageParts {
	title, body, has := strings.Cut(self.String(), "\n")
	if !has {
		body = ""
	}
	for strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return CommitMessageParts{
		Subject: title,
		Text:    body,
	}
}

// String implements the fmt.Stringer interface.
func (self CommitMessage) String() string {
	return string(self)
}
