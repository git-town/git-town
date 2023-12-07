package validate

import (
	"errors"

	"github.com/git-town/git-town/v11/src/messages"
)

// HasGitVersion verifies that the system has Git of version 2.7 or newer installed.
func HasGitVersion(majorVersion, minorVersion int) error {
	if !IsAcceptableGitVersion(majorVersion, minorVersion) {
		return errors.New(messages.GitVersionTooLow)
	}
	return nil
}

// IsAcceptableGitVersion indicates whether the given Git version works for Git Town.
func IsAcceptableGitVersion(major, minor int) bool {
	return major > 2 || (major == 2 && minor >= 7)
}
