package git

import "github.com/git-town/git-town/v9/src/gohacks/stringslice"

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config          RepoConfig
	Backend         BackendCommands
	Frontend        FrontendCommands
	CommandsCounter Counter
	FinalMessages   stringslice.Collector
}

type Counter interface {
	Count() int
}
