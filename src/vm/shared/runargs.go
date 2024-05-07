package shared

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
)

type RunArgs struct {
	Backend                         git.BackendCommands
	Config                          *config.ValidatedConfig
	Connector                       hostingdomain.Connector
	DialogTestInputs                components.TestInputs
	FinalMessages                   stringslice.Collector
	Frontend                        git.FrontendCommands
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	UpdateInitialBranchLocalSHA     func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
