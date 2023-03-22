// Package git runs Git commands in a controlled and typesafe way.
package git

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/subshell"
)

func NewInternalRunner(dir *string, debug bool) InternalRunner {
	internalRunnerInstance := subshell.InternalRunner{Dir: dir}
	if debug {
		return subshell.InternalDebuggingRunner{InternalRunner: internalRunnerInstance}
	}
	return internalRunnerInstance
}

// NewPublicCommands provides a PublicRepo instance that behaves according to the given configuration and uses the given InternalRepo instance internally.
func NewPublicRunner(omitBranchNames, dryRun bool, currentBranchCache *cache.String) PublicRunner {
	if dryRun {
		return subshell.PublicDryRunner{
			CurrentBranch:   currentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	}
	return subshell.PublicRunner{
		CurrentBranch:   currentBranchCache,
		OmitBranchNames: omitBranchNames,
	}
}
