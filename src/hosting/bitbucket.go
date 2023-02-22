package hosting

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/giturl"
)

// BitbucketConnector provides access to the API of Bitbucket installations.
type BitbucketConnector struct {
	Config
	organization string
	git          gitRunner
}

// NewBitbucketConnector provides a Bitbucket driver instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewBitbucketConnector(url giturl.Parts, gitConfig gitConfig, git gitRunner) *BitbucketConnector {
	manualOrigin := gitConfig.OriginOverride()
	if manualOrigin != "" {
		url.Host = manualOrigin
	}
	if gitConfig.HostingService() != "bitbucket" && url.Host != "bitbucket.org" {
		return nil
	}
	return &BitbucketConnector{
		Config: Config{
			apiToken:   "",
			hostname:   url.Host,
			originURL:  gitConfig.OriginURL(),
			owner:      url.Org,
			repository: url.Repo,
		},
		organization: url.Org,
		git:          git,
	}
}

func (c *BitbucketConnector) ProposalForBranch(branch string) (*Proposal, error) {
	return nil, fmt.Errorf("BitBucket API functionality isn't implemented yet")
}

func (c *BitbucketConnector) DefaultCommitMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (c *BitbucketConnector) HostingServiceName() string {
	return "Bitbucket"
}

func (c *BitbucketConnector) NewProposalURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	branchSha, err := c.git.ShaForBranch(branch)
	if err != nil {
		return "", fmt.Errorf("cannot determine pull request URL from %q to %q: %w", branch, parentBranch, err)
	}
	query.Add("source", strings.Join([]string{c.organization + "/" + c.repository, branchSha[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{c.organization + "/" + c.repository, "", parentBranch}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", c.RepositoryURL(), query.Encode()), nil
}

func (c *BitbucketConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.organization, c.repository)
}

//nolint:nonamedreturns
func (c *BitbucketConnector) SquashMergeProposal(number int, message string) (mergeSHA string, err error) {
	return "", errors.New("shipping pull requests via the Bitbucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (c *BitbucketConnector) UpdateProposalTarget(number int, target string) error {
	return errors.New("shipping pull requests via the Bitbucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}
