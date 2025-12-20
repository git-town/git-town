package shared

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type RunArgs struct {
	Backend                         subshelldomain.RunnerQuerier
	BranchInfos                     gitdomain.BranchInfos
	Config                          Mutable[config.ValidatedConfig]
	Connector                       Option[forgedomain.Connector]
	FinalMessages                   stringslice.Collector
	Frontend                        subshelldomain.Runner
	Git                             git.Commands
	Inputs                          dialogcomponents.Inputs
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	UpdateInitialSnapshotLocalSHA   func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
