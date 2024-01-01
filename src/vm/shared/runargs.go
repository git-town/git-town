package shared

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
)

type RunArgs struct {
	PrependOpcodes                  func(...Opcode)
	Connector                       hostingdomain.Connector
	Lineage                         configdomain.Lineage
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	Runner                          *git.ProdRunner
	UpdateInitialBranchLocalSHA     func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
