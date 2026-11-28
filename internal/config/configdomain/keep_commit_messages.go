package configdomain

// KeepCommitMessages indicates whether to preserve the commit messages of squashed commits.
type KeepCommitMessages bool

// Enabled indicates whether keeping commit messages is enabled.
func (self KeepCommitMessages) Enabled() bool {
	return bool(self)
}
