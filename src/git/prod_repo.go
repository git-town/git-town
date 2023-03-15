package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/stringslice"
	"github.com/git-town/git-town/v7/src/subshell"
)

// ProdRepo is a Git Repo in production code.
type ProdRepo struct {
	Config        *config.GitTown // the git.Configuration instance for this repo
	DryRun        *subshell.DryRun
	Logging       Repo                    // the Runner instance to Git operations that show up in the output
	LoggingRunner *subshell.LoggingRunner // the LoggingRunner instance used
	Silent        Repo                    // the Runner instance for silent Git operations
}

// NewProdRepo provides a Repo instance in the current working directory.
func NewProdRepo(debugFlag *bool) ProdRepo {
	silentRunner := subshell.InternalDebuggingRunner{Debug: debugFlag}
	config := config.NewGitTown(silentRunner)
	currentBranchTracker := cache.String{}
	dryRun := subshell.DryRun{}
	isRepoCache := cache.Bool{}
	remoteBranchCache := cache.Strings{}
	remotesCache := cache.Strings{}
	silentRepo := Repo{
		Runner:             silentRunner,
		Config:             &config,
		CurrentBranchCache: &currentBranchTracker,
		DryRun:             &dryRun,
		IsRepoCache:        &isRepoCache,
		RemotesCache:       &remotesCache,
		RemoteBranchCache:  &remoteBranchCache,
		RootDirCache:       &cache.String{},
	}
	loggingRunner := subshell.NewLoggingRunner(&silentRepo, &dryRun)
	loggingRepo := Repo{
		Runner:             loggingRunner,
		Config:             &config,
		CurrentBranchCache: &currentBranchTracker,
		DryRun:             &dryRun,
		IsRepoCache:        &isRepoCache,
		RemotesCache:       &remotesCache,
		RemoteBranchCache:  &remoteBranchCache,
		RootDirCache:       &cache.String{},
	}
	return ProdRepo{
		Silent:        silentRepo,
		Logging:       loggingRepo,
		LoggingRunner: loggingRunner,
		Config:        &config,
		DryRun:        &dryRun,
	}
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (r *ProdRepo) RemoveOutdatedConfiguration() error {
	branches, err := r.Silent.LocalAndOriginBranches()
	if err != nil {
		return err
	}
	for child, parent := range r.Config.ParentBranchMap() {
		hasChildBranch := stringslice.Contains(branches, child)
		hasParentBranch := stringslice.Contains(branches, parent)
		if !hasChildBranch || !hasParentBranch {
			err = r.Config.RemoveParentBranch(child)
			if err != nil {
				return err
			}
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
