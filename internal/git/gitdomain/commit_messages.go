package gitdomain

import "strings"

type CommitMessages []CommitMessage

func NewCommitMessages(messages ...string) CommitMessages {
	result := make(CommitMessages, len(messages))
	for m, message := range messages {
		result[m] = CommitMessage(message)
	}
	return result
}

func (self CommitMessages) Join() CommitMessage {
	var result strings.Builder
	for c, message := range self {
		if c > 0 {
			result.WriteByte('\n')
		}
		result.WriteString(message.String())
		if c < len(self)-1 && !strings.HasSuffix(message.String(), "\n") {
			result.WriteByte('\n')
		}
	}
	return CommitMessage(result.String())
}
