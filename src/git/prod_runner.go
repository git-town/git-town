package git

import "github.com/git-town/git-town/v9/src/gohacks"

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config          RepoConfig
	Backend         BackendCommands
	Frontend        FrontendCommands
	CommandsCounter *gohacks.Counter
}
