package git

import (
	"github.com/git-town/git-town/v10/src/gohacks"
	"github.com/git-town/git-town/v10/src/gohacks/stringslice"
)

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config          RepoConfig
	Backend         BackendCommands
	Frontend        FrontendCommands
	CommandsCounter *gohacks.Counter
	FinalMessages   *stringslice.Collector
}
