package shared

import (
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type RunArgs struct {
	Backend                         gitdomain.RunnerQuerier
	Config                          config.ValidatedConfig
	Connector                       Option[hostingdomain.Connector]
	DialogTestInputs                components.TestInputs
	FinalMessages                   stringslice.Collector
	Frontend                        gitdomain.Runner
	Git                             git.Commands
	InitialBranchesSnapshot         Option[gitdomain.BranchesSnapshot]
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	UpdateInitialBranchLocalSHA     func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
