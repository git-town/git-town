package shared

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

type RunArgs struct {
	PrependOpcodes                  func(...Opcode)
	Connector                       hosting.Connector
	Lineage                         config.Lineage
	RegisterUndoablePerennialCommit func(domain.SHA)
	Runner                          *git.ProdRunner
	UpdateInitialBranchLocalSHA     func(domain.LocalBranchName, domain.SHA) error
}
