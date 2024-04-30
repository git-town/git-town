package git

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
)

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Backend         BackendCommands
	CommandsCounter *gohacks.Counter
	Config          *configdomain.ValidatedConfig
	FinalMessages   *stringslice.Collector
	Frontend        FrontendCommands
}
