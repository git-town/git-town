package shared

import (
	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
)

type RunArgs struct {
	Connector                       hostingdomain.Connector
	DialogTestInputs                *components.TestInputs
	Lineage                         configdomain.Lineage
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	Runner                          *git.ProdRunner
	UpdateInitialBranchLocalSHA     func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
