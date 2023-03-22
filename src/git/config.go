package git

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
)

// RepoConfig represents the known state of a Git repository.
type RepoConfig struct {
	*config.GitTown
	CurrentBranchCache *cache.String  // caches the currently checked out Git branch
	DryRun             bool           // single source of truth for whether to dry-run Git commands in this repo
	IsRepoCache        *cache.Bool    // caches whether the current directory is a Git repo
	RemoteBranchCache  *cache.Strings // caches the remote branches of this Git repo
	RemotesCache       *cache.Strings // caches Git remotes
	RootDirCache       *cache.String  // caches the base of the Git directory
}

func NewRepoConfig(runner InternalRunner, dryRun bool) RepoConfig {
	return RepoConfig{
		GitTown:            config.NewGitTown(runner),
		CurrentBranchCache: &cache.String{},
		DryRun:             dryRun,
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
}
