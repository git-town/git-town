package commands

// CreateStandaloneTag creates a tag not on a branch.
func CreateStandaloneTag(repo Repo, name string) error {
	return repo.RunMany([][]string{
		{"git", "checkout", "-b", "temp"},
		{"touch", "a.txt"},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "temp"},
		{"git", "tag", "-a", name, "-m", name},
		{"git", "checkout", "-"},
		{"git", "branch", "-D", "temp"},
	})
}
