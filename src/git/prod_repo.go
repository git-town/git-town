package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/run"
)

// ProdRepo is a Git Repo in production code.
type ProdRepo struct {
	Config       config.Config // the git.Configuration instance for this repo
	DryRun       *DryRun
	Logging      Runner        // the Runner instance to Git operations that show up in the output
	LoggingShell *LoggingShell // the LoggingShell instance used
	Silent       Runner        // the Runner instance for silent Git operations
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
	for child, parent := range r.Config.ParentBranchMap() {
		hasChildBranch, err := r.Silent.HasLocalOrOriginBranch(child)
		if err != nil {
			return err
		}
		hasParentBranch, err := r.Silent.HasLocalOrOriginBranch(parent)
		if err != nil {
			return err
		}
		if !hasChildBranch || !hasParentBranch {
			return r.Config.DeleteParentBranch(child)
		}
	}
	return nil
}

// NavigateToRootIfNecessary changes into the root directory of the current repository.
func (r *ProdRepo) NavigateToRootIfNecessary() error {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current working directory: %w", err)
	}
	gitRootDirectory, err := r.Silent.RootDirectory()
	if err != nil {
		return err
	}
	if currentDirectory == gitRootDirectory {
		return nil
	}
	return os.Chdir(gitRootDirectory)
}
