package validate

import (
	"errors"

	"github.com/git-town/git-town/v8/src/git"
)

// HasGitVersion verifies that the system has Git of version 2.7 or newer installed.
func HasGitVersion(backend *git.BackendCommands) error {
	majorVersion, minorVersion, err := backend.Version()
	if err != nil {
		return err
	}
	if !IsAcceptableGitVersion(majorVersion, minorVersion) {
		return errors.New("this app requires Git 2.7.0 or higher")
	}
	return nil
}

// IsAcceptableGitVersion indicates whether the given Git version works for Git Town.
func IsAcceptableGitVersion(major, minor int) bool {
	return major > 2 || (major == 2 && minor >= 7)
}
