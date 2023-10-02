package domain

// RepoRootDir represents the root directory of a Git repository.
type RepoRootDir struct {
	value string
}

func NewRepoRootDir(dir string) RepoRootDir {
	return RepoRootDir{value: dir}
}

func (rrd RepoRootDir) IsEmpty() bool {
	return rrd.value == ""
}

func (rrd RepoRootDir) String() string {
	return rrd.value
}
