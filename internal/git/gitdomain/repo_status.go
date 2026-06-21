package gitdomain

type RepoStatus struct {
	Conflicts        bool // the repo contains merge conflicts
	OpenChanges      bool // there are uncommitted changes
	RebaseInProgress bool // a rebase is in progress
	UntrackedChanges bool // the repo contains files that aren't tracked by Git
}

// NeedsToCommit indicates whether the repo contains something to commit.
func (self RepoStatus) NeedsToCommit() bool {
	return self.OpenChanges || self.RebaseInProgress || self.Conflicts || self.UntrackedChanges
}
