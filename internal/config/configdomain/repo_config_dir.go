package configdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// RepoConfigDir is the directory that contains the Git Town configuration for a particular Git repo.
// Example: ~/.config/git-town/home-user-my-repo.
type RepoConfigDir stringss.TrimmedString

func (self RepoConfigDir) String() string {
	return string(self)
}
