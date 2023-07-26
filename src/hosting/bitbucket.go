package hosting

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/messages"
)

// BitbucketConnector provides access to the API of Bitbucket installations.
type BitbucketConnector struct {
	CommonConfig
	git gitCommands
}

// NewBitbucketConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewBitbucketConnector(gitConfig gitTownConfig, git gitCommands) (*BitbucketConnector, error) {
	hostingService, err := gitConfig.HostingService()
	if err != nil {
		return nil, err
	}
	url := gitConfig.OriginURL()
	if url == nil || (url.Host != "bitbucket.org" && hostingService != config.HostingServiceBitbucket) {
		return nil, nil //nolint:nilnil
	}
	return &BitbucketConnector{
		CommonConfig: CommonConfig{
			APIToken:     "",
			Hostname:     url.Host,
			Organization: url.Org,
			Repository:   url.Repo,
		},
		git: git,
	}, nil
}

func (c *BitbucketConnector) FindProposal(_, _ string) (*Proposal, error) {
	return nil, fmt.Errorf(messages.HostingBitBucketNotImplemented)
}

func (c *BitbucketConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (c *BitbucketConnector) HostingServiceName() string {
	return "Bitbucket"
}

func (c *BitbucketConnector) NewProposalURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	branchSha, err := c.git.ShaForBranch(branch)
	if err != nil {
		return "", fmt.Errorf(messages.ProposalURLProblem, branch, parentBranch, err)
	}
	query.Add("source", strings.Join([]string{c.Organization + "/" + c.Repository, branchSha[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{c.Organization + "/" + c.Repository, "", parentBranch}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", c.RepositoryURL(), query.Encode()), nil
}

func (c *BitbucketConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.Hostname, c.Organization, c.Repository)
}

func (c *BitbucketConnector) SquashMergeProposal(_ int, _ string) (mergeSHA string, err error) {
	return "", errors.New(messages.HostingBitBucketNotImplemented)
}

func (c *BitbucketConnector) UpdateProposalTarget(_ int, _ string) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}
