// Package git runs Git commands in a controlled and typesafe way.
package git

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/subshell"
)

// NewInternalRepo provides InternalRepo instances that optionally debugs depending on the given parameter.
func NewInternalRepo(debug bool) InternalRepo {
	shellRunner := subshell.InternalRunner{Dir: nil}
	var gitRunner InternalRunner
	if debug {
		gitRunner = subshell.InternalDebuggingRunner{InternalRunner: shellRunner}
	} else {
		gitRunner = shellRunner
	}
	return InternalRepo{
		InternalRunner:     gitRunner,
		Config:             config.NewGitTown(gitRunner),
		CurrentBranchCache: &cache.String{},
		DryRun:             false,
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
}

// NewPublicRepo provides a PublicRepo instance that behaves according to the given configuration and uses the given InternalRepo instance internally.
func NewPublicRepo(omitBranchNames, dryRun bool, internalRepo InternalRepo) PublicRepo {
	var gitRunner PublicRunner
	if dryRun {
		gitRunner = subshell.PublicDryRunner{
			CurrentBranch:   internalRepo.CurrentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	} else {
		gitRunner = subshell.PublicRunner{
			CurrentBranch:   internalRepo.CurrentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	}
	return PublicRepo{
		Public:       gitRunner,
		InternalRepo: internalRepo,
	}
}
