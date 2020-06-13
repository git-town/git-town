package git

import "github.com/git-town/git-town/src/command"

// ProdRepo is a Git Repo in production code.
type ProdRepo struct {
	Silent         Runner        // the Runner instance for silent Git operations
	Logging        Runner        // the Runner instance to Git operations that show up in the output
	LoggingShell   *LoggingShell // the LoggingShell instance used
	*Configuration               // the git.Configuration instance for this repo
}

// NewProdRepo provides a Repo instance in the current working directory.
func NewProdRepo() *ProdRepo {
	silentShell := command.SilentShell{}
	config := Config()
	currentBranchTracker := CurrentBranchTracker{}
	remoteBranchCache := RemoteBranchCache{}
	remotesCache := RemotesCache{}
	silentRunner := NewRunner(silentShell, config, &currentBranchTracker, &remotesCache, &remoteBranchCache)
	loggingShell := NewLoggingShell()
	loggingRunner := NewRunner(loggingShell, config, &currentBranchTracker, &remotesCache, &remoteBranchCache)
	return &ProdRepo{
		Silent:        silentRunner,
		Logging:       loggingRunner,
		LoggingShell:  loggingShell,
		Configuration: config,
	}
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (r *ProdRepo) RemoveOutdatedConfiguration() error {
	for child, parent := range r.GetParentBranchMap() {
		hasChildBranch, err := r.Silent.HasLocalOrRemoteBranch(child)
		if err != nil {
			return err
		}
		hasParentBranch, err := r.Silent.HasLocalOrRemoteBranch(parent)
		if err != nil {
			return err
		}
		if !hasChildBranch || !hasParentBranch {
			r.Configuration.DeleteParentBranch(child)
		}
	}
	return nil
}
