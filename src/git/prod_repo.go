package git

import "github.com/git-town/git-town/src/command"

// ProdRepo is a Git Repo in production code.
type ProdRepo struct {
	Silent Runner // the Runner instance for silent commands
	// loggingRunner  Runner // the Runner instance to logging commands
	*Configuration // the git.Configuration instance for this repo
}

// NewProdRepo provides a Repo instance in the current working directory.
func NewProdRepo() *ProdRepo {
	shell := command.SilentShell{}
	config := NewConfiguration(shell)
	silentRunner := Runner{
		Shell:          shell,
		currentBranch:  &CurrentBranchTracker{},
		remoteBranches: &RemoteBranchTracker{},
		Configuration:  config,
	}
	return &ProdRepo{Silent: silentRunner, Configuration: config}
}
