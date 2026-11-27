package gitdomain

import (
	"strings"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
)

type CommitMessages []CommitMessage

func NewCommitMessages(messages ...string) CommitMessages {
	result := make(CommitMessages, len(messages))
	for m, message := range messages {
		result[m] = CommitMessage(message)
	}
	return result
}

func (self CommitMessages) Join(sep string) CommitMessage {
	messages := make([]string, len(self))
	for c, message := range self {
		messages[c] = message.String()
	}
	return CommitMessage(stringss.Trim(strings.Join(messages, sep)))
}
