// Package git runs Git commands in a controlled and typesafe way.
package git

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/subshell"
)

func NewBackendRunner(dir *string, debug bool) BackendRunner {
	backendRunner := subshell.BackendRunner{Dir: dir}
	if debug {
		return subshell.BackendLoggingRunner{Runner: backendRunner}
	}
	return backendRunner
}

// NewFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func NewFrontendRunner(omitBranchNames, dryRun bool, currentBranchCache *cache.String) FrontendRunner {
	if dryRun {
		return subshell.FrontendDryRunner{
			CurrentBranch:   currentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	}
	return subshell.FrontendRunner{
		CurrentBranch:   currentBranchCache,
		OmitBranchNames: omitBranchNames,
	}
}
