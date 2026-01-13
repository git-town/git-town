package configdomain

// RepoConfigDir is the directory that contains the Git Town configuration for a particular Git repo.
// Example: ~/.config/git-town/home-user-my-repo.
type RepoConfigDir string

func (self RepoConfigDir) String() string {
	return string(self)
}
