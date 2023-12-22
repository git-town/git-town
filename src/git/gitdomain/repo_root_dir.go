package gitdomain

// RepoRootDir represents the root directory of a Git repository.
type RepoRootDir string

func EmptyRepoRootDir() RepoRootDir {
	return ""
}

func NewRepoRootDir(dir string) RepoRootDir {
	return RepoRootDir(dir)
}

func (self RepoRootDir) IsEmpty() bool {
	return self == ""
}

func (self RepoRootDir) String() string {
	return string(self)
}
