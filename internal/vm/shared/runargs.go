package shared

import (
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type RunArgs struct {
	Backend                         subshelldomain.RunnerQuerier
	BranchInfos                     Option[gitdomain.BranchInfos]
	Config                          Mutable[config.ValidatedConfig]
	Connector                       Option[forgedomain.Connector]
	Detached                        configdomain.Detached
	FinalMessages                   stringslice.Collector
	Frontend                        subshelldomain.Runner
	Git                             git.Commands
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	UpdateInitialSnapshotLocalSHA   func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
