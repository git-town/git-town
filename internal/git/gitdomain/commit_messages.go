package gitdomain

type CommitMessages []CommitMessage

func NewCommitMessages(messages ...string) CommitMessages {
	result := make(CommitMessages, len(messages))
	for m, message := range messages {
		result[m] = CommitMessage(message)
	}
	return result
}
