package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/config"
	"github.com/git-town/git-town/src/run"
)

// ProdRepo is a Git Repo in production code.
type ProdRepo struct {
	Silent       Runner         // the Runner instance for silent Git operations
	Logging      Runner         // the Runner instance to Git operations that show up in the output
	LoggingShell *LoggingShell  // the LoggingShell instance used
	Config       *config.Config // the git.Configuration instance for this repo
	DryRun       *DryRun
}

// NewProdRepo provides a Repo instance in the current working directory.
func NewProdRepo() *ProdRepo {
	silentShell := run.SilentShell{}
	config := config.NewConfiguration(silentShell)
	currentBranchTracker := StringCache{}
	dryRun := DryRun{}
	isRepoCache := BoolCache{}
	remoteBranchCache := StringSliceCache{}
	remotesCache := StringSliceCache{}
	silentRunner := Runner{
		Shell:              silentShell,
		Config:             config,
		CurrentBranchCache: &currentBranchTracker,
		DryRun:             &dryRun,
		IsRepoCache:        &isRepoCache,
		RemotesCache:       &remotesCache,
		RemoteBranchCache:  &remoteBranchCache,
		RootDirCache:       &StringCache{},
	}
	loggingShell := NewLoggingShell(&silentRunner, &dryRun)
	loggingRunner := Runner{
		Shell:              loggingShell,
		Config:             config,
		CurrentBranchCache: &currentBranchTracker,
		DryRun:             &dryRun,
		IsRepoCache:        &isRepoCache,
		RemotesCache:       &remotesCache,
		RemoteBranchCache:  &remoteBranchCache,
		RootDirCache:       &StringCache{},
	}
	return &ProdRepo{
		Silent:       silentRunner,
		Logging:      loggingRunner,
		LoggingShell: loggingShell,
		Config:       config,
		DryRun:       &dryRun,
	}
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (r *ProdRepo) RemoveOutdatedConfiguration() error {
	for child, parent := range r.Config.GetParentBranchMap() {
		hasChildBranch, err := r.Silent.HasLocalOrRemoteBranch(child)
		if err != nil {
			return err
		}
		hasParentBranch, err := r.Silent.HasLocalOrRemoteBranch(parent)
		if err != nil {
			return err
		}
		if !hasChildBranch || !hasParentBranch {
			return r.Config.DeleteParentBranch(child)
		}
	}
	return nil
}

func (r *ProdRepo) NavigateToRootIfNecessary() error {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current working directory: %w", err)
	}
	gitRootDirectory, err := r.Silent.RootDirectory()
	if err != nil {
		return err
	}
	if currentDirectory != gitRootDirectory {
		return os.Chdir(gitRootDirectory)
	}
	return nil
}
