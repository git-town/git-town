package git

import (
	"github.com/git-town/git-town/src/command"
)

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
	silentRunner := Runner{
		Shell:             silentShell,
		currentBranch:     &currentBranchTracker,
		remoteBranchCache: &RemoteBranchCache{},
		Configuration:     config,
	}
	loggingShell := NewLoggingShell(&currentBranchTracker)
	loggingRunner := Runner{
		Shell:         loggingShell,
		currentBranch: &CurrentBranchTracker{},
		Configuration: config,
	}
	return &ProdRepo{Silent: silentRunner, Logging: loggingRunner, LoggingShell: loggingShell, Configuration: config}
}
