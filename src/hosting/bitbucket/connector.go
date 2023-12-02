package bitbucket

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git/giturl"
	"github.com/git-town/git-town/v10/src/hosting/common"
	"github.com/git-town/git-town/v10/src/messages"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	common.Config
	getSHAForBranch common.SHAForBranchFunc
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) (*Connector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "bitbucket.org" && args.HostingService != config.HostingBitbucket) {
		return nil, nil //nolint:nilnil
	}
	return &Connector{
		Config: common.Config{
			APIToken:     "",
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		getSHAForBranch: args.GetSHAForBranch,
	}, nil
}

type NewConnectorArgs struct {
	OriginURL       *giturl.Parts
	HostingService  config.Hosting
	GetSHAForBranch common.SHAForBranchFunc
}

func (self *Connector) DefaultProposalMessage(proposal domain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self *Connector) FindProposal(_, _ domain.LocalBranchName) (*domain.Proposal, error) {
	return nil, fmt.Errorf(messages.HostingBitBucketNotImplemented)
}

func (self *Connector) HostingServiceName() string {
	return "Bitbucket"
}

func (self *Connector) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	return fmt.Sprintf("%s/pull-requests/new?source=%s&dest=%s%%2F%s%%3A%s",
			self.RepositoryURL(),
			url.QueryEscape(branch.String()),
			url.QueryEscape(self.Organization),
			url.QueryEscape(self.Repository),
			url.QueryEscape(parentBranch.String())),
		nil
}

func (self *Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.Hostname, self.Organization, self.Repository)
}

func (self *Connector) SquashMergeProposal(_ int, _ string) (mergeSHA domain.SHA, err error) {
	return domain.EmptySHA(), errors.New(messages.HostingBitBucketNotImplemented)
}

func (self *Connector) UpdateProposalTarget(_ int, _ domain.LocalBranchName) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}
