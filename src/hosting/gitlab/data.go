package gitlab

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
)

type Data struct {
	hostingdomain.Data
	APIToken Option[configdomain.GitLabToken]
}

func (self Data) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (!%d)", proposal.Title, proposal.Number)
}

func (self Data) NewProposalURL(branch, parentBranch, _ gitdomain.LocalBranchName, _, _ string) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch.String())
	query.Add("merge_request[target_branch]", parentBranch.String())
	return fmt.Sprintf("%s/-/merge_requests/new?%s", self.RepositoryURL(), query.Encode()), nil
}

func (self Data) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", self.baseURL(), self.projectPath())
}

func (self Data) baseURL() string {
	return "https://" + self.HostnameWithStandardPort()
}

func (self Data) projectPath() string {
	return fmt.Sprintf("%s/%s", self.Organization, self.Repository)
}
