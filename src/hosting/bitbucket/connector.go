package bitbucket

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	hostingdomain.Config
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) (*Connector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "bitbucket.org" && args.HostingService != configdomain.HostingPlatformBitbucket) {
		return nil, nil //nolint:nilnil
	}
	return &Connector{
		Config: hostingdomain.Config{
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
	}, nil
}

type NewConnectorArgs struct {
	OriginURL      *giturl.Parts
	HostingService configdomain.HostingPlatform
}

func (self *Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self *Connector) FindProposal(_, _ gitdomain.LocalBranchName) (*hostingdomain.Proposal, error) {
	return nil, fmt.Errorf(messages.HostingBitBucketNotImplemented)
}

func (self *Connector) HostingServiceName() string {
	return "Bitbucket"
}

func (self *Connector) NewProposalURL(branch, parentBranch gitdomain.LocalBranchName) (string, error) {
	return fmt.Sprintf("%s/pull-requests/new?source=%s&dest=%s%%2F%s%%3A%s",
			self.RepositoryURL(),
			url.QueryEscape(branch.String()),
			url.QueryEscape(self.Organization),
			url.QueryEscape(self.Repository),
			url.QueryEscape(parentBranch.String())),
		nil
}

func (self *Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self *Connector) SquashMergeProposal(_ int, _ string) (mergeSHA gitdomain.SHA, err error) {
	return gitdomain.EmptySHA(), errors.New(messages.HostingBitBucketNotImplemented)
}

func (self *Connector) UpdateProposalTarget(_ int, _ gitdomain.LocalBranchName) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}
