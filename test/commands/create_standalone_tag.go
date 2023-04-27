package commands

import "github.com/git-town/git-town/v8/test/subshell"

// CreateStandaloneTag creates a tag not on a branch.
func CreateStandaloneTag(shell subshell.Mocking, name string) error {
	return shell.RunMany([][]string{
		{"git", "checkout", "-b", "temp"},
		{"touch", "a.txt"},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "temp"},
		{"git", "tag", "-a", name, "-m", name},
		{"git", "checkout", "-"},
		{"git", "branch", "-D", "temp"},
	})
}
