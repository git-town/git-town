package bitbucketcloud

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

var (
	bbclWebConnector AnonConnector
	_                forgedomain.Connector = bbclWebConnector
)

// AnonConnector provides connectivity to Bitbucket Cloud without authentication.
type AnonConnector struct {
	forgedomain.Data
}

func (self AnonConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner)
	return nil
}

func (self AnonConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self AnonConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pull-requests/new?source=%s&dest=%s%%2F%s%%3A%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(self.Organization),
		url.QueryEscape(self.Repository),
		url.QueryEscape(data.ParentBranch.String()))
}

func (self AnonConnector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self AnonConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
