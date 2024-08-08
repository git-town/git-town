package bitbucket

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	hostingdomain.Data
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	return Connector{
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
}

type NewConnectorArgs struct {
	HostingPlatform Option[configdomain.HostingPlatform]
	RemoteURL       giturl.Parts
}

func (self Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self Connector) FindProposal(_, _ gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	return None[hostingdomain.Proposal](), errors.New(messages.HostingBitBucketNotImplemented)
}

func (self Connector) NewProposalURL(branch, parentBranch, _ gitdomain.LocalBranchName, _ gitdomain.ProposalTitle, _ gitdomain.ProposalBody) (string, error) {
	return fmt.Sprintf("%s/pull-requests/new?source=%s&dest=%s%%2F%s%%3A%s",
			self.RepositoryURL(),
			url.QueryEscape(branch.String()),
			url.QueryEscape(self.Organization),
			url.QueryEscape(self.Repository),
			url.QueryEscape(parentBranch.String())),
		nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SquashMergeProposal(_ int, _ gitdomain.CommitMessage) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}

func (self Connector) UpdateProposalTarget(_ int, _ gitdomain.LocalBranchName) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}
