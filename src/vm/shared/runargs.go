package shared

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/hosting"
)

type RunArgs struct {
	PrependOpcodes                  func(...Opcode)
	Connector                       hosting.Connector
	Lineage                         configdomain.Lineage
	RegisterUndoablePerennialCommit func(domain.SHA)
	Runner                          *git.ProdRunner
	UpdateInitialBranchLocalSHA     func(domain.LocalBranchName, domain.SHA) error
}
