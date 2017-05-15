package git

import "github.com/Originate/git-town/src/util"

// EnsureIsRepository asserts that the current directory is in a repository
func EnsureIsRepository() {
	util.Ensure(IsRepository(), "This is not a Git repository.")
}

// IsRepository returns whether or not the current directory is in a repository
func IsRepository() bool {
	_, err := util.GetFullCommandOutput("git", "rev-parse")
	return err == nil
}
