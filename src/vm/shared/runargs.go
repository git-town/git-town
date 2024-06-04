package shared

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
)

type RunArgs struct {
	Backend                         gitdomain.RunnerQuerier
	Config                          config.ValidatedConfig
	Connector                       Option[hostingdomain.Connector]
	DialogTestInputs                components.TestInputs
	FinalMessages                   stringslice.Collector
	Frontend                        gitdomain.Runner
	Git                             git.Commands
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	UpdateInitialBranchLocalSHA     func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
