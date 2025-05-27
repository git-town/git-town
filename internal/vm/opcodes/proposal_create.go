package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// ProposalCreate creates a new proposal for the current branch.
type ProposalCreate struct {
	Branch                  gitdomain.LocalBranchName
	MainBranch              gitdomain.LocalBranchName
	ProposalBody            gitdomain.ProposalBody
	ProposalTitle           gitdomain.ProposalTitle
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalCreate) Run(args shared.RunArgs) error {
	parentBranch, hasParentBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParentBranch {
		return fmt.Errorf(messages.ProposalNoParent, self.Branch)
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	prURL, err := connector.NewProposalURL(self.Branch, parentBranch, self.MainBranch, self.ProposalTitle, self.ProposalBody)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Frontend, args.Backend)
	return nil
}
