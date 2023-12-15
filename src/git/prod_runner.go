package git

import (
	"github.com/git-town/git-town/v11/src/config"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/gohacks/stringslice"
)

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	*config.GitTown
	Backend         BackendCommands
	Frontend        FrontendCommands
	CommandsCounter *gohacks.Counter
	FinalMessages   *stringslice.Collector
}
