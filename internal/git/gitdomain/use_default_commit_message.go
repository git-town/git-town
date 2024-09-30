package gitdomain

// FallbackToDefaultCommitMessage defines whether to fall back to the default commit message when no commit message is given.
type FallbackToDefaultCommitMessage bool

const (
	FallbackToDefaultCommitMessageNo  FallbackToDefaultCommitMessage = false // if the commit message is missing, let the user enter it
	FallbackToDefaultCommitMessageYes FallbackToDefaultCommitMessage = true  // if the commit message is missing, use the default commit message
)
