package gitdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// RepoRootDir represents the root directory of a Git repository.
type RepoRootDir stringss.Trimmed

func NewRepoRootDir(dir string) RepoRootDir {
	if dir == "" {
		panic("empty root dir provided")
	}
	return RepoRootDir(dir)
}

func (self RepoRootDir) String() string {
	return string(self)
}
