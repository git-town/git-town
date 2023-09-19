package domain

// RepoRootDir represents the root directory of a Git repository.
type RepoRootDir struct {
	value string
}

func NewRepoRootDir(dir string) RepoRootDir {
	return RepoRootDir{value: dir}
}

func (r RepoRootDir) IsEmpty() bool {
	return r.value == ""
}

func (r RepoRootDir) String() string {
	return r.value
}
