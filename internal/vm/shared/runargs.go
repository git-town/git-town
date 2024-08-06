package shared

import (
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/config"
	"github.com/git-town/git-town/v14/internal/git"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	. "github.com/git-town/git-town/v14/internal/gohacks/prelude"
	"github.com/git-town/git-town/v14/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v14/internal/hosting/hostingdomain"
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
