package git

import "github.com/git-town/git-town/v9/src/git"

// A zero-value SHA to be used as a placeholder in tests.
func ZeroValueSHA() git.SHA {
	return git.ErrorSHA()
}
