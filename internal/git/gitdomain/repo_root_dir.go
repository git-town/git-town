package gitdomain

// RepoRootDir represents the root directory of a Git repository.
type RepoRootDir string

func NewRepoRootDir(dir string) RepoRootDir {
	if dir == "" {
		panic("empty root dir provided")
	}
	return RepoRootDir(dir)
}

func (self RepoRootDir) String() string {
	return string(self)
}
