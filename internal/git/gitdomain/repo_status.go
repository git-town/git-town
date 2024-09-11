package gitdomain

type RepoStatus struct {
	Conflicts        bool // the repo contains merge conflicts
	OpenChanges      bool // there are uncommitted changes
	RebaseInProgress bool // a rebase is in progress
}
