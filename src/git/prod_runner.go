package git

import (
	"github.com/git-town/git-town/v12/src/config"
	"github.com/git-town/git-town/v12/src/gohacks"
	"github.com/git-town/git-town/v12/src/gohacks/stringslice"
)

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	*config.Config
	Backend         BackendCommands
	Frontend        FrontendCommands
	CommandsCounter *gohacks.Counter
	FinalMessages   *stringslice.Collector
}
