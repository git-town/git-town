package git

import "github.com/Originate/git-town/lib/util"

// EnsureIsRepository asserts that the current directory is in a repository
func EnsureIsRepository() {
	util.Ensure(isRepository(), "This is not a git repository.")
}

// Helpers

func isRepository() bool {
	_, err := util.GetFullCommandOutput("git", "rev-parse")
	return err == nil
}
