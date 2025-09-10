package gitea

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

var giteaAnonConnector AnonConnector
var _ forgedomain.Connector = giteaAnonConnector

// AnonConnector provides connectivity to gitea without authentication.
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
	toCompare := data.ParentBranch.String() + "..." + data.Branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare))
}

func (self AnonConnector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self AnonConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
