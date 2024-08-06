package gitdomain

type Commits []Commit

// ContainsSHA indicates whether this commits list contains a commit with the given SHA.
func (self Commits) ContainsSHA(sha SHA) bool {
	for _, commit := range self {
		if commit.SHA == sha {
			return true
		}
	}
	return false
}

func (self Commits) Messages() CommitMessages {
	result := make(CommitMessages, len(self))
	for c, commit := range self {
		result[c] = commit.Message
	}
	return result
}

func (self Commits) SHAs() SHAs {
	result := make(SHAs, len(self))
	for c, commit := range self {
		result[c] = commit.SHA
	}
	return result
}
