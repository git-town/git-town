package gitdomain

import "strings"

type CommitMessage string

// CommitMessageParts describes the parts of a Git commit message.
type CommitMessageParts struct {
	Subject string
	Text    string
}

// Parts separates the parts of the given commit message.
func (self CommitMessage) Parts() CommitMessageParts {
	parts := strings.SplitN(self.String(), "\n", 2)
	title := parts[0]
	body := ""
	if len(parts) == 2 {
		body = parts[1]
	}
	for strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return CommitMessageParts{
		Subject: title,
		Text:    body,
	}
}

func (self CommitMessage) String() string {
	return string(self)
}
