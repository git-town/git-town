package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/browser"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
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
		return hostingdomain.UnsupportedServiceError()
	}
	prURL, err := connector.NewProposalURL(self.Branch, parentBranch, self.MainBranch, self.ProposalTitle, self.ProposalBody)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Frontend, args.Backend)
	return nil
}
