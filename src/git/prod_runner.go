package git

import "github.com/git-town/git-town/v7/src/subshell"

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
	Stats    *subshell.Statistics
}
