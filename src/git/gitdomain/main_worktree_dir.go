package gitdomain

// MainWorkTreeDir represents the root directory of the main worktree of this Git repo.
type MainWorkTreeDir string

func EmptyMainWorkTreeDir() MainWorkTreeDir {
	return ""
}

func NewMainWorkTreeDir(dir string) MainWorkTreeDir {
	return MainWorkTreeDir(dir)
}

func (self MainWorkTreeDir) IsEmpty() bool {
	return self == ""
}

func (self MainWorkTreeDir) String() string {
	return string(self)
}
