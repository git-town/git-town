package domain

// RepoRootDir represents the root directory of a Git repository.
type RepoRootDir struct {
	value string
}

func EmptyRepoRootDir() RepoRootDir {
	return NewRepoRootDir("")
}

func NewRepoRootDir(dir string) RepoRootDir {
	return RepoRootDir{value: dir}
}

func (self RepoRootDir) IsEmpty() bool {
	return self.value == ""
}

func (self RepoRootDir) String() string {
	return self.value
}
