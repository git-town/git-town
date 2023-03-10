package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/git-town/git-town/v7/src/stringslice"
)

// ProdRepo is a Git Repo in production code.
type ProdRepo struct {
	Config       config.GitTown // the git.Configuration instance for this repo
	DryRun       *run.DryRun
	Logging      Repo              // the Repo instance to Git operations that show up in the output
	LoggingShell *run.LoggingShell // the LoggingShell instance used
	Silent       Repo              // the Repo instance for silent Git operations
}

// NewProdRepo provides a Repo instance in the current working directory.
func NewProdRepo(debugFlag *bool) ProdRepo {
	silentShell := run.SilentShell{Debug: debugFlag}
	config := config.NewGitTown(silentShell)
	currentBranchTracker := cache.String{}
	dryRun := run.DryRun{}
	isRepoCache := cache.Bool{}
	remoteBranchCache := cache.Strings{}
	remotesCache := cache.Strings{}
	silentRepo := Repo{
		Shell:              silentShell,
		Config:             config,
		CurrentBranchCache: &currentBranchTracker,
		DryRun:             &dryRun,
		IsRepoCache:        &isRepoCache,
		RemotesCache:       &remotesCache,
		RemoteBranchCache:  &remoteBranchCache,
		RootDirCache:       &cache.String{},
	}
	loggingShell := run.NewLoggingShell(&silentRepo, &dryRun)
	loggingRepo := Repo{
		Shell:              loggingShell,
		Config:             config,
		CurrentBranchCache: &currentBranchTracker,
		DryRun:             &dryRun,
		IsRepoCache:        &isRepoCache,
		RemotesCache:       &remotesCache,
		RemoteBranchCache:  &remoteBranchCache,
		RootDirCache:       &cache.String{},
	}
	return ProdRepo{
		Silent:       silentRepo,
		Logging:      loggingRepo,
		LoggingShell: loggingShell,
		Config:       config,
		DryRun:       &dryRun,
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
