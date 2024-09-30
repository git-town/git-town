package gitdomain

// UseDefaultCommitMessage defines possible ways to deal with a missing commit message.
type UseDefaultCommitMessage bool

const (
	UseDefaultCommitMessageNo  UseDefaultCommitMessage = false // if the commit message is missing, let the user enter it
	UseDefaultCommitMessageYes UseDefaultCommitMessage = true  // if the commit message is missing, use the default commit message
)
