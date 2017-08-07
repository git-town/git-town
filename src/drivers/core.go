package drivers

import "github.com/Originate/git-town/src/git"

// Core provides the public API for the drivers subsystem

var registry = &Registry{}
var activeDriver *CodeHostingDriver

// GetActiveDriver returns the code hosting driver to use
func GetActiveDriver() *CodeHostingDriver {
	if activeDriver == nil {
		activeDriver = registry.DetermineActiveDriver(git.GetURLHostname(git.GetRemoteOriginURL()))
	}
	return activeDriver
}
