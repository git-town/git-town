package bitbucketdatacenter

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

// type-check to ensure conformance to the Connector interface
var (
	bbdcWebConnector WebConnector
	_                forgedomain.Connector = bbdcWebConnector
)

// WebConnector provides connectivity to Bitbucket DataCenter without authentication.
type WebConnector struct {
	forgedomain.Data
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner)
	return nil
}

func (self WebConnector) DefaultProposalMessage(proposalData forgedomain.ProposalData) string {
	data := proposalData.Data()
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pull-requests?create&sourceBranch=%s&targetBranch=%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(data.ParentBranch.String()))
}

func (self WebConnector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/projects/%s/repos/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
