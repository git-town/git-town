package gitlab

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
)

type Config struct {
	hostingdomain.Config
	APIToken configdomain.GitLabToken
}

func (self *Config) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (!%d)", proposal.Title, proposal.Number)
}

func (self *Config) NewProposalURL(branch, parentBranch gitdomain.LocalBranchName) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch.String())
	query.Add("merge_request[target_branch]", parentBranch.String())
	return fmt.Sprintf("%s/-/merge_requests/new?%s", self.RepositoryURL(), query.Encode()), nil
}

func (self *Config) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", self.baseURL(), self.projectPath())
}

func (self *Config) baseURL() string {
	return "https://" + self.HostnameWithStandardPort()
}

func (self *Config) projectPath() string {
	return fmt.Sprintf("%s/%s", self.Organization, self.Repository)
}
